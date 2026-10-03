import 'package:finance_app/model/check.dart';
import 'package:finance_app/model/product.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/json.dart';

import '../model/purchase.dart';
import '../model/stats.dart';

class PurchaseService {
  static Future<PurchaseModel> getPurchase(int id) {
    return Server.get('/purchase/status', query: {'id': id}).then((resp) {
      return PurchaseModel.fromJson(jsonDec(resp));
    });
  }

  static Future<List<CheckResultModel>> checkDestination(String checkerId,
      List<ProductDestinationFieldRequestModel> destination) async {
    Map<String, String> query = {"id": checkerId};
    for (int i = 0; i < destination.length; i++) {
      query["input.$i.key"] = destination[i].key;
      query["input.$i.value"] = destination[i].value;
    }
    var resp = await Server.get('/checker/', query: query);
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => CheckResultModel.fromJson(e))
        .toList();
  }

  static Future<int> createPurchase(
      int productId, List<ProductDestinationFieldRequestModel> destination,
      {String? pin}) async {
    Map<String, dynamic> body = {"productId": productId.toString()};
    var i = 0;
    for (var dst in destination) {
      body["destination.$i.key"] = dst.key;
      body["destination.$i.value"] = dst.value;
      i++;
    }
    if (pin != null) {
      body["pin"] = pin;
    }
    var resp = await Server.post('/purchase/create', body: body);
    return jsonDec(resp)["purchaseId"];
  }

  static Future<PurchaseBillModel> checkBill(int productId,
      List<ProductDestinationFieldRequestModel> destination) async {
    Map<String, dynamic> body = {
      "productId": productId.toString(),
      "pay": "false"
    };
    var i = 0;
    for (var dst in destination) {
      body["destination.$i.key"] = dst.key;
      body["destination.$i.value"] = dst.value;
      i++;
    }
    var resp = await Server.post('/purchase/bill', body: body);
    return PurchaseBillModel.fromJson(jsonDec(resp));
  }

  static Future<int> payBill(
      int productId,
      int billId,
      List<ProductDestinationFieldRequestModel> destination,
      String? pin) async {
    Map<String, dynamic> body = {
      "productId": productId.toString(),
      "pay": "true",
      "id": billId.toString(),
      "pin": pin ?? ""
    };
    var i = 0;
    for (var dst in destination) {
      body["destination.$i.key"] = dst.key;
      body["destination.$i.value"] = dst.value;
      i++;
    }
    var resp = await Server.post('/purchase/bill', body: body);
    return jsonDec(resp)["purchaseId"];
  }

  static Future<List<PurchaseModel>> getPurchases(
      {int? status, int? limit, int? beforeId}) async {
    Map<String, dynamic> query = {};
    if (status != null) {
      query["status"] = status;
    }
    if (limit != null) {
      query["limit"] = limit;
    }
    if (beforeId != null) {
      query["beforeId"] = beforeId;
    }
    var resp = await Server.get('/purchase/history', query: query);
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => PurchaseModel.fromJson(e))
        .toList();
  }

  static Future<String> updateUserSellPrice(int id, int sellPrice) async {
    var resp = await Server.post('/purchase/updateUserSellPrice', body: {
      "id": id.toString(),
      "userSellPrice": sellPrice.toString(),
    });
    return jsonDec(resp)["message"];
  }

  static Future<List<PurchaseModel>> searchPurchases(String query) async {
    var resp = await Server.get('/purchase/search', query: {"query": query});
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => PurchaseModel.fromJson(e))
        .toList();
  }

  static Future<PurchaseStatsWrapperModel> getStats() async {
    var resp = await Server.get('/purchase/stats');
    return PurchaseStatsWrapperModel.fromJson(jsonDec(resp));
  }
}
