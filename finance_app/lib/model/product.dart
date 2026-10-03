import 'package:flutter/cupertino.dart';

class ProductModel {
  int id;
  int categoryId;
  int type;
  String categoryName;
  String name;
  String kind;
  String description;
  int price;
  int? billAdmin;
  String imageUrl;
  bool isAvailable;
  int destinationId;
  int? beforeDiscountPrice;

  ProductModel({
    required this.id,
    required this.categoryName,
    required this.categoryId,
    required this.type,
    required this.name,
    required this.kind,
    required this.description,
    required this.price,
    required this.billAdmin,
    required this.imageUrl,
    required this.isAvailable,
    required this.destinationId,
    required this.beforeDiscountPrice,
  });

  static List<ProductModel> sortByPrice(List<ProductModel> products,
      {bool lowest = true}) {
    List<ProductModel> products2 = List.from(products);
    products2.sort((a, b) {
      if (!lowest) {
        return b.price.compareTo(a.price);
      }
      return a.price.compareTo(b.price);
    });
    return products2;
  }

  static List<ProductModel> sortByNameAZ(List<ProductModel> products) {
    List<ProductModel> products2 = List.from(products);
    products2.sort((a, b) {
      return a.name.compareTo(b.name);
    });
    return products2;
  }

  factory ProductModel.fromJson(Map<String, dynamic> json) {
    return ProductModel(
      id: json['id'],
      categoryName: json['categoryName'],
      categoryId: json['categoryId'],
      type: json['type'],
      name: json['name'],
      kind: json['kind'],
      description: json['description'],
      price: json['price'],
      billAdmin: json['billAdmin'],
      imageUrl: json['imageUrl'],
      isAvailable: json['isAvailable'],
      destinationId: json['destinationId'],
      beforeDiscountPrice: json['beforeDiscountPrice'],
    );
  }

  bool hasImage() {
    return imageUrl != "";
  }

  bool isPostpaid() {
    return type == 2;
  }
}

class ProductCategoryModel {
  int id;
  int? parentId;
  String name;
  String description;
  String iconUrl;
  int productType;
  List<ProductCategoryModel> children;
  Map<String, String> meta;

  ProductCategoryModel({
    required this.id,
    this.parentId,
    required this.name,
    required this.description,
    required this.iconUrl,
    required this.productType,
    required this.children,
    required this.meta,
  });

  bool isPostpaid() {
    return productType == 2;
  }

  bool lastDepth() {
    for (var child in children) {
      if (child.children.isNotEmpty) {
        return false;
      }
    }
    return true;
  }

  bool hasIcon() {
    return iconUrl != "";
  }

  String getMeta(String key) {
    return meta[key] ?? "";
  }

  bool isGrid() {
    return getMeta("layout") == "grid";
  }

  factory ProductCategoryModel.fromJson(Map<String, dynamic> json) {
    return ProductCategoryModel(
        id: json['id'],
        name: json['name'],
        description: json['description'],
        iconUrl: json['iconUrl'],
        productType: json['productType'],
        meta: (json['meta'] as Map<String, dynamic>)
            .map((key, value) => MapEntry(key, value.toString())),
        children: ((json['children'] ?? []) as List<dynamic>)
            .map((category) => ProductCategoryModel.fromJson(category))
            .toList());
  }
}

class ProductDestinationFieldOptionModel {
  String label;
  String value;

  ProductDestinationFieldOptionModel(
      {required this.label, required this.value});

  factory ProductDestinationFieldOptionModel.fromJson(
      Map<String, dynamic> json) {
    return ProductDestinationFieldOptionModel(
        label: json['label'], value: json['value']);
  }
}

class ProductDestinationFieldModel {
  static const String textType = "text";
  static const String optionType = "option";
  static const String phoneType = "phone";
  static const String numberType = "number";

  String key;
  String label;
  String description;
  String type;
  List<ProductDestinationFieldOptionModel>? options;
  bool required;

  ProductDestinationFieldModel(
      {required this.key,
      required this.label,
      required this.description,
      required this.type,
      this.options,
      required this.required});

  factory ProductDestinationFieldModel.fromJson(Map<String, dynamic> json) {
    return ProductDestinationFieldModel(
        key: json['key'],
        label: json['label'],
        description: json['description'],
        type: json['type'],
        options: ((json['options'] ?? []) as List<dynamic>)
            .map((options) =>
                ProductDestinationFieldOptionModel.fromJson(options))
            .toList(),
        required: json['required']);
  }

  TextInputType getKeyboardType() {
    switch (type) {
      case numberType:
      case phoneType:
        return TextInputType.number;
    }
    return TextInputType.text;
  }

  bool isOption() {
    return type == optionType;
  }

  bool isPhone() {
    return type == phoneType;
  }
}

class ProductDestinationModel {
  int id;
  String name;
  String description;
  String? checkerId;
  List<ProductDestinationFieldModel> fields;

  ProductDestinationModel(
      {required this.id,
      required this.name,
      required this.description,
      this.checkerId,
      required this.fields});

  factory ProductDestinationModel.fromJson(Map<String, dynamic> json) {
    return ProductDestinationModel(
        id: json['id'],
        name: json['name'],
        description: json['description'],
        checkerId: json['checkerId'],
        fields: ((json['fields'] ?? []) as List<dynamic>)
            .map((field) => ProductDestinationFieldModel.fromJson(field))
            .toList());
  }
}

class ProductDestinationFieldRequestModel {
  String key;
  String? label;
  String? labelValue;
  String value;

  ProductDestinationFieldRequestModel(this.key, this.value,
      {this.label, this.labelValue});
}
