import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:flutter/material.dart';

class Alert {
  static Future<bool?> showYesOrNo({
    required BuildContext context,
    required String title,
    required String message,
    VoidCallback? onYes,
    VoidCallback? onNo,
  }) {
    return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (context) => AlertDialog(
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
        ),
        insetPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 0),
        title: Text(title,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Text(
          message,
          style: const TextStyle(fontSize: 14),
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context, false);
              if (onNo != null) {
                onNo();
              }
            },
            child: const Text("Tidak"),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context, true);
              if (onYes != null) {
                onYes();
              }
            },
            child: const Text("Ya"),
          ),
        ],
      ),
    );
  }

  static Future<T?> show<T>({
    bool barrierDismissible = true,
    required String title,
    required String message,
    String? closeMessage,
    List<Widget>? actions,
  }) {
    return showDialog<T>(
      context: Screens.navigatorKey.currentContext!,
      barrierDismissible: barrierDismissible,
      builder: (context) => AlertDialog(
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
        ),
        title: Text(title,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Text(message, style: const TextStyle(fontSize: 14)),
        actions: actions ??
            [
              TextButton(
                onPressed: () {
                  Navigator.pop(context);
                },
                child: Text(closeMessage ?? "Tutup"),
              ),
            ],
      ),
    );
  }

  static bool isLoadingShow = false;

  static void closeLoading({BuildContext? context}) {
    if (isLoadingShow) {
      Navigator.of(context ?? Screens.navigatorKey.currentContext!).pop();
    }
    isLoadingShow = false;
  }

  static Future<void> withLoading(Function(Function() done) func,
      {String? message}) async {
    Alert.showLoading(message: message);
    try {
      await func(() => Alert.closeLoading());
    } catch (e) {
      Alert.closeLoading();
      Alert.message(e.toString());
      rethrow;
    }
  }

  static Future<T?> error<T>(Exception e) async {
    return message(e.toString());
  }

  static Future<T?> message<T>(String message,
      {String btnMsg = "Tutup"}) async {
    return showDialog<T>(
        context: Screens.navigatorKey.currentContext!,
        builder: (context) {
          return AlertDialog(
            shape: const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(Radius.circular(8)),
            ),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const SizedBox(height: 16),
                Padding(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                    child: Text(message,
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                            fontSize: 14.5, fontWeight: FontWeight.w600))),
                const SizedBox(height: 16),
                Button(
                    onPressed: () {
                      Navigator.pop(context, true);
                    },
                    width: double.infinity,
                    child: Text(btnMsg)),
                const SizedBox(height: 12),
              ],
            ),
          );
        });
  }

  static void showLoading({String? message, BuildContext? context}) {
    if (isLoadingShow) {
      Screens.back();
    }
    isLoadingShow = true;
    showDialog(
        context: context ?? Screens.navigatorKey.currentContext!,
        barrierDismissible: false,
        builder: (context) => WillPopScope(
            onWillPop: () async => false,
            child: Dialog(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16),
              ),
              child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 32),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const CircularProgressIndicator(),
                    const SizedBox(height: 24),
                    Text(message ?? "Loading..."),
                  ],
                ),
              ),
            ))).whenComplete(() => isLoadingShow = false);
  }
}
