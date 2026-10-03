import 'package:finance_app/screens/balance/deposit_info_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/utils/alert.dart';

import '../utils/screens.dart';

class NotificationModel {
  static const int typeUnknown = 0;
  static const int typeSystem = 1;
  static const int typePurchase = 2;
  static const int typeDeposit = 3;
  static const int typeTransfer = 4;
  static const int typeWithdraw = 5;

  int id;
  int type;
  String title;
  String description;
  bool hasRead;
  Map<String, String> data;
  int createdAt;

  NotificationModel(
      {required this.id,
      required this.type,
      required this.title,
      required this.description,
      required this.data,
      required this.hasRead,
      required this.createdAt});

  factory NotificationModel.fromJson(Map<String, dynamic> json) =>
      NotificationModel(
          id: json["id"],
          type: json["type"],
          title: json["title"],
          description: json["description"],
          hasRead: json["hasRead"],
          data: (json["data"] as Map<String, dynamic>)
              .map((key, value) => MapEntry(key, value.toString())),
          createdAt: json["createdAt"]);

  void onClick() {
    switch (type) {
      case NotificationModel.typeDeposit:
        Alert.withLoading((done) async {
          final deposit = await BalanceService.getDepositDetailed(
              int.parse(data["depositId"]!));
          done();
          Screens.to(DepositInfoScreen(deposit: deposit));
        });
        break;
      default:
        break;
    }
  }
}
