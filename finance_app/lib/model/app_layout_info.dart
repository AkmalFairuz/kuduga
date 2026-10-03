import 'package:finance_app/model/product.dart';
import 'package:finance_app/screens/product/product_category_list_screen.dart';
import 'package:finance_app/screens/product/product_list_screen.dart';
import 'package:finance_app/screens/product/product_postpaid_list_screen.dart';
import 'package:finance_app/service/product_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';

class AppLayoutInfo {
  List<ServiceInfoWrapper> services;

  AppLayoutInfo({required this.services});

  factory AppLayoutInfo.fromJson(Map<String, dynamic> json) {
    return AppLayoutInfo(
        services: (json['services'] as List<dynamic>)
            .map((e) => ServiceInfoWrapper.fromJson(e))
            .toList());
  }
}

class AppVersionInfo {
  int version;
  int minimumVersion;

  AppVersionInfo({required this.version, required this.minimumVersion});

  factory AppVersionInfo.fromJson(Map<String, dynamic> json) {
    return AppVersionInfo(
        version: json['version'], minimumVersion: json['minimumVersion']);
  }
}

class AppBannerInfo {
  String image;
  String? url;

  AppBannerInfo({required this.image, this.url});

  factory AppBannerInfo.fromJson(Map<String, dynamic> json) {
    return AppBannerInfo(image: json['image'], url: json['url']);
  }
}

class ServiceInfoWrapper {
  String title;
  List<ServiceInfo> services;

  ServiceInfoWrapper({required this.services, required this.title});

  factory ServiceInfoWrapper.fromJson(Map<String, dynamic> json) {
    return ServiceInfoWrapper(
        title: json['title'],
        services: (json['services'] as List<dynamic>)
            .map((e) => ServiceInfo.fromJson(e))
            .toList());
  }
}

class ServiceInfo {
  String actionType;
  Map<String, String> actionData;
  String name;
  String iconType;
  String iconData;
  dynamic _actionCache;

  ServiceInfo(
      {required this.actionType,
      required this.actionData,
      required this.iconType,
      required this.iconData,
      required this.name});

  factory ServiceInfo.fromJson(Map<String, dynamic> json) {
    return ServiceInfo(
      actionType: json['actionType'],
      actionData: (json['actionData'] as Map<String, dynamic>)
          .map((key, value) => MapEntry(key, value.toString())),
      iconType: json['iconType'],
      iconData: json['iconData'],
      name: json['name'],
    );
  }

  Future<void> doAction() async {
    switch (actionType) {
      case "pulsa":
        int? pulsaParentId;
        if (actionData["parentId"] != null) {
          pulsaParentId = int.parse(actionData["parentId"]!);
        }
        Screens.to(ProductListScreen(
            title: name,
            category: null,
            isPulsa: true,
            parentPulsaCategoryId: pulsaParentId));
        return;
      case "product":
        ProductCategoryModel? category;
        if (_actionCache == null) {
          await Alert.withLoading((done) async {
            _actionCache = await ProductService.getProductCategory(
                int.parse(actionData["id"]!));
            category = _actionCache;
            done();
          });
        } else {
          category = _actionCache;
        }
        if (category!.isPostpaid()) {
          Screens.to(ProductPostpaidListScreen(category: category!));
        } else {
          Screens.to(ProductListScreen(
              title: name, category: category, isPulsa: false));
        }
        return;
      case "category":
        var categories = <ProductCategoryModel>[];
        if (_actionCache == null) {
          await Alert.withLoading((done) async {
            categories = await ProductService.getProductCategories(
                nested: true, parentId: int.parse(actionData["id"]!));
            _actionCache = categories;
            done();
          });
        } else {
          categories = _actionCache;
        }
        Screens.to(ProductCategoryListScreen(
            title: name,
            isGrid: actionData["isGrid"] == "true",
            categories: categories));
        return;
    }
  }
}
