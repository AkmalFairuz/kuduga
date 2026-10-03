import 'package:finance_app/model/deposit.dart';
import 'package:finance_app/model/payment_method.dart';
import 'package:finance_app/model/transaction.dart';
import 'package:finance_app/model/transfer.dart';

import '../server/server.dart';
import '../utils/json.dart';

class BalanceService {
  static Future<int> getBalance() {
    return Server.get('/balance').then((resp) {
      if (resp.statusCode != 200) {
        throw Exception('Failed to get balance: ${resp.body}');
      }
      return jsonDec(resp) as int;
    });
  }

  static Future<DepositDetailedModel> createDeposit(
      int amount, int paymentMethod) {
    return Server.post('/balance/deposit/create', body: {
      'amount': amount.toString(),
      'paymentMethod': paymentMethod.toString()
    }).then((resp) {
      if (resp.statusCode != 200) {
        throw Exception('Failed to deposit: ${resp.body}');
      }
      return DepositDetailedModel.fromJson(jsonDec(resp));
    });
  }

  static Future<List<DepositModel>> getDepositHistory({int? status}) {
    Map<String, dynamic> query = {};
    if (status != null) {
      query['filterStatus'] = status.toString();
    }
    return Server.get('/balance/deposit/history', query: query).then((resp) {
      if (resp.statusCode != 200) {
        throw Exception('Failed to get deposit history: ${resp.body}');
      }
      List<dynamic> json = jsonDec(resp);
      return json.map((e) => DepositModel.fromJson(e)).toList();
    });
  }

  static Future<void> cancelDeposit(int id) {
    return Server.post('/balance/deposit/cancel', query: {'id': id.toString()})
        .then((resp) {
      if (resp.statusCode != 200) {
        throw Exception('Failed to cancel deposit: ${resp.body}');
      }
    });
  }

  static Future<List<TransactionModel>> getTransactions(
      {required int from, required int to, int? limit = 300}) {
    return Server.get('/balance/transactions', query: {
      'from': from.toString(),
      'to': to.toString(),
      'limit': limit
    }).then((resp) {
      if (resp.statusCode != 200) {
        throw Exception('Failed to get transaction history: ${resp.body}');
      }
      List<dynamic> json = jsonDec(resp);
      return json.map((e) => TransactionModel.fromJson(e)).toList();
    });
  }

  static Future<DepositDetailedModel> getDepositDetailed(int id) {
    return Server.get('/balance/deposit/detail', query: {'id': id})
        .then((resp) {
      return DepositDetailedModel.fromJson(jsonDec(resp));
    });
  }

  static Future<List<PaymentMethodModel>> getDepositPaymentMethods() async {
    var resp = await Server.get("/balance/deposit/paymentMethods");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => PaymentMethodModel.fromJson(e))
        .toList();
  }

  static Future<List<TransferUserInfoModel>> getTransferDestinations() async {
    var resp = await Server.get("/balance/transfer/destinations");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => TransferUserInfoModel.fromJson(e))
        .toList();
  }

  static Future<TransferUserInfoModel> getTransferDestination(
      String email) async {
    var resp = await Server.get("/balance/transfer/checkDestination",
        query: {"email": email});
    return TransferUserInfoModel.fromJson(jsonDec(resp));
  }

  static Future<TransferModel> transferToUser({
    required TransferUserInfoModel destination,
    required int amount,
    required String note,
    required int purpose,
    String? pin,
  }) async {
    var body = {
      "destinationEmail": destination.email,
      "amount": amount.toString(),
      "note": note,
      "purpose": purpose.toString()
    };
    if (pin != null) {
      body["pin"] = pin;
    }
    var resp = await Server.post("/balance/transfer", body: body);
    return TransferModel.fromJson(jsonDec(resp));
  }

  static Future<List<TransferModel>> getTransfers() async {
    var resp = await Server.get("/balance/transfer/history");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => TransferModel.fromJson(e))
        .toList();
  }
}
