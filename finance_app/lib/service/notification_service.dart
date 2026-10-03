import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/model/notification.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/utils/json.dart';
import 'package:finance_app/utils/native_bindings.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

typedef MessageHandler = void Function(String, Map<String, dynamic>);

class NotificationService {
  static Map<String, String> notificationChannels = {
    "other": "Lainnya",
    "promotions": "Promosi",
    "announcements": "Pengumuman",
    "transactions": "Transaksi",
  };

  static Future<List<NotificationModel>> getNotifications() async {
    var resp = await Server.get("/notification/");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => NotificationModel.fromJson(e))
        .toList();
  }

  static NotificationService instance = NotificationService();

  String? _fcmToken;
  final FlutterLocalNotificationsPlugin _notificationsPlugin =
      FlutterLocalNotificationsPlugin();
  int _notificationId = 1;
  final List<MessageHandler> _messageHandlers = [];

  NotificationService();

  void init() {
    notificationChannels.forEach((channelId, channelName) {
      NativeBindings.channel.invokeMethod(
          "createNotificationChannel", {"id": channelId, "name": channelName});
    });

    requestPermissions();
    registerListener();

    const InitializationSettings initializationSettingsAndroid =
        InitializationSettings(
            android: AndroidInitializationSettings("@mipmap/ic_launcher"),
            iOS: DarwinInitializationSettings());
    _notificationsPlugin.initialize(
      initializationSettingsAndroid,
      onDidReceiveNotificationResponse: (details) {
        if (details.input != null) {}
      },
    );
  }

  Future<void> showLocal(
      {String? title, String? body, String? imageUrl, String? channelId}) {
    String channelId2 = channelId ?? "other";
    String? channelName = notificationChannels[channelId2];
    if (channelName == null) {
      channelId2 = "other";
      channelName = "Lainnya";
    }

    return _notificationsPlugin.show(
        _notificationId++, // TODO: implement notification ID
        title,
        body,
        NotificationDetails(
            android: AndroidNotificationDetails(channelId2, channelName,
                importance: Importance.high,
                priority: Priority.high,
                // TODO: implement image
                visibility: NotificationVisibility.public)));
  }

  Future<bool> requestPermissions() async {
    NotificationSettings settings =
        await FirebaseMessaging.instance.requestPermission();
    return settings.authorizationStatus == AuthorizationStatus.authorized;
  }

  void registerListener() {
    FirebaseMessaging.instance.onTokenRefresh.listen((fcmToken) {
      if (Server.hasToken()) {
        pushToken(fcmToken);
      }
    });

    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      var data = message.data;
      if (data.containsKey("_event")) {
        handleMessage(data["_event"], data);
      }
      var notification = message.notification;
      if (notification != null) {
        showLocal(title: notification.title, body: notification.body);
      }
    });
  }

  void handleMessage(String event, Map<String, dynamic> data) {
    for (final handler in _messageHandlers) {
      handler(event, data);
    }
    StateNotifier.notify("notification");
    switch (event) {
      case "refreshBalance":
        AccountController.getInstance().fetchBalance();
        break;
      case "deposit":
        AccountController.getInstance().fetchBalance();
        StateNotifier.notify("deposit");
        break;
      case "purchaseStatusUpdate":
        StateNotifier.notify("purchase");
        break;
    }
  }

  void subscribeMessage(MessageHandler callback) {
    _messageHandlers.add(callback);
  }

  void unsubscribeMessage(MessageHandler callback) {
    _messageHandlers.remove(callback);
  }

  Future<void> pushToken([String? token]) async {
    token ??= await FirebaseMessaging.instance.getToken();
    _fcmToken = token;

    await Server.post("/notification/updateFcmToken",
        body: {"fcmToken": _fcmToken});
  }
}
