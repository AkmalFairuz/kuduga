import 'dart:async';

import 'package:finance_app/model/product.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/json.dart';

class ProductService {
  static Future<List<ProductModel>> getProductByCategory(int categoryId) {
    return Server.get("/product", query: {"categoryId": categoryId})
        .then((res) {
      List<dynamic> decoded = jsonDec(res);
      return decoded.map((x) => ProductModel.fromJson(x)).toList();
    });
  }

  static Future<List<ProductCategoryModel>> getProductCategories(
      {bool nested = false, int? parentId}) {
    return Server.get("/product/categories",
        query: {"isNested": nested, "parentId": parentId}).then((res) {
      List<dynamic> decoded = jsonDec(res);
      return decoded.map((x) => ProductCategoryModel.fromJson(x)).toList();
    });
  }

  static Future<Map<String, String>> getPulsaCategories() async {
    final resp = await Server.get("/product/pulsaCategories");
    return (jsonDec(resp) as Map<String, dynamic>)
        .map((key, value) => MapEntry(key, value.toString()));
  }

  static Future<ProductCategoryModel> getProductCategory(int categoryId) async {
    var resp = await Server.get("/product/category", query: {"id": categoryId});
    return ProductCategoryModel.fromJson(jsonDec(resp));
  }

  static Future<ProductDestinationModel> getProductDestination(int id) async {
    var resp = await Server.get("/product/destination", query: {"id": id});
    return ProductDestinationModel.fromJson(jsonDec(resp));
  }

  static Future<ProductModel> getProduct(int id) async {
    var resp = await Server.get("/product/single", query: {"id": id});
    return ProductModel.fromJson(jsonDec(resp));
  }
}
