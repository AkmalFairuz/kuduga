import 'package:finance_app/service/printer/monospace_formatter.dart';
import 'package:finance_app/service/printer/printable.dart';
import 'package:finance_app/service/printer/thermal_printer_formatter.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class PurchaseModel implements ThermalPrintable {
  static const int statusWaitingPayment = 0;
  static const int statusProcess = 1;
  static const int statusSuccess = 2;
  static const int statusFailed = 3;

  int purchaseId;
  int productId;
  int productType;
  String productSku;
  String productName;
  String productDescription;
  int productCategoryId;
  String productCategoryName;
  String productKind;
  int status;
  String statusText;
  Map<String, String> destination;
  Map<String, dynamic> extraData;
  String note;

  List<Map<String, String>>? billData;
  int? totalBillFee;
  int? totalBill;
  int? billAdmin;

  String? proof;
  int price;
  int userSellPrice;
  int fee;
  int? paymentMethod;
  int paymentExpiredAt;
  String paymentData;
  int createdAt;

  PurchaseModel({
    required this.purchaseId,
    required this.productId,
    required this.productSku,
    required this.productName,
    required this.productType,
    required this.productDescription,
    required this.productCategoryId,
    required this.productCategoryName,
    required this.productKind,
    required this.status,
    required this.statusText,
    required this.destination,
    required this.extraData,
    required this.note,
    this.billData,
    this.totalBill,
    this.totalBillFee,
    this.billAdmin,
    this.proof,
    required this.price,
    required this.userSellPrice,
    required this.fee,
    this.paymentMethod,
    required this.paymentExpiredAt,
    required this.paymentData,
    required this.createdAt,
  });

  static Map<String, String> _toDestination(Map<String, dynamic> json) {
    Map<String, String> resultMap = {};
    json.forEach((key, value) {
      resultMap[key] = value.toString();
    });
    return resultMap;
  }

  factory PurchaseModel.fromJson(Map<String, dynamic> json) {
    return PurchaseModel(
      purchaseId: json['purchaseId'],
      productId: json['productId'],
      productSku: json['productSku'],
      productName: json['productName'],
      productType: json['productType'],
      productDescription: json['productDescription'],
      productCategoryId: json['productCategoryId'],
      productCategoryName: json['productCategoryName'],
      productKind: json['productKind'],
      status: json['status'],
      statusText: json['statusText'],
      destination: _toDestination(json['destination']),
      billData: json['billData'] != null
          ? (json['billData'] as List<dynamic>)
              .map((e) => (e as Map<String, dynamic>)
                  .map((key, value) => MapEntry(key, value.toString())))
              .toList()
          : null,
      totalBill: json['totalBill'],
      totalBillFee: json['totalBillFee'],
      billAdmin: json['billAdmin'],
      extraData: json['extraData'] as Map<String, dynamic>,
      note: json['note'],
      proof: json['proof'],
      price: json['price'],
      userSellPrice: json['userSellPrice'],
      fee: json['fee'],
      paymentMethod: json['paymentMethod'],
      paymentExpiredAt: json['paymentExpiredAt'],
      paymentData: json['paymentData'],
      createdAt: json['createdAt'],
    );
  }

  bool hasProof() {
    return proof != null && proof != "";
  }

  bool isBill() {
    return productType == 2;
  }

  Color statusColor() {
    switch (status) {
      case PurchaseModel.statusWaitingPayment:
      case PurchaseModel.statusProcess:
        return Colors.amber[900]!;
      case PurchaseModel.statusSuccess:
        return Meta.color;
      case PurchaseModel.statusFailed:
        return Colors.red[600]!;
    }
    return Colors.grey[600]!;
  }

  IconData statusIcon() {
    switch (status) {
      case PurchaseModel.statusWaitingPayment:
      case PurchaseModel.statusProcess:
        return Icons.autorenew;
      case PurchaseModel.statusSuccess:
        return Icons.check;
      case PurchaseModel.statusFailed:
        return Icons.priority_high;
    }
    return Icons.question_mark;
  }

  @override
  String thermalPrint(int maxCharacterPerLine) {
    final wrapper = ThermalPrinterFormatter();
    wrapper.setMaxCharactersPerLine(maxCharacterPerLine);

    final f1 = MonospaceFormatter(
        nameLength: 10, valueLength: maxCharacterPerLine - 13);
    f1.add("PRODUK", productName);
    f1.add("KATEGORI",
        "$productCategoryName${productKind.isNotEmpty ? " - $productKind" : ""}");
    if (isBill()) {
      f1.add("TAGIHAN", Meta.currencyFormatRp(totalBill! + billAdmin!));
      f1.add("ADM LOKET", Meta.currencyFormatRp(userSellPrice));
      f1.add("TOTAL",
          Meta.currencyFormatRp(totalBill! + billAdmin! + userSellPrice));
    } else {
      f1.add("HARGA", Meta.currencyFormatRp(userSellPrice));
    }
    for (final dst in destination.entries) {
      f1.add(dst.key.toUpperCase(), dst.value);
    }
    for (final data in extraData.entries) {
      if (data.key.toUpperCase() == "JUMLAH TAGIHAN") {
        continue;
      }
      f1.add(data.key.toUpperCase(), data.value.toString());
    }
    if (isBill() && billData!.length == 1) {
      for (final data in billData!.first.entries) {
        if (data.key.toUpperCase() == "NILAI TAGIHAN") {
          continue;
        }
        f1.add(data.key.toUpperCase(), data.value.toString());
      }
    }
    wrapper.append(
        "#$purchaseId   ${Meta.formatUnixDate(createdAt, format: "d MMM yyyy HH:mm")}");
    wrapper.line();
    wrapper.stripes();
    wrapper.line();
    wrapper.rawAppend(f1.format());
    wrapper.line();

    if (isBill() && billData!.length > 1) {
      var billNo = 0;
      for (final bd in billData!) {
        billNo++;
        wrapper.stripes();
        wrapper.line();
        wrapper.rawAppend("[C]<b>Detail Tagihan #$billNo</b>");
        wrapper.line();
        final f2 = MonospaceFormatter(
            nameLength: 10, valueLength: maxCharacterPerLine - 13);
        for (final data in bd.entries) {
          f2.add(data.key.toUpperCase(), data.value.toString());
        }
        wrapper.rawAppend(f2.format());
        wrapper.line();
      }
    }

    wrapper.stripes();
    wrapper.line();
    if (proof != null && proof!.isNotEmpty) {
      wrapper.line();
      wrapper.rawAppend("[C]<b>* Serial Number *</b>");
      wrapper.line();
      final trimmedProof = proof!.trim();
      if (trimmedProof.length > maxCharacterPerLine) {
        wrapper.appendFormatted(trimmedProof, prefix: "[C]", suffix: "");
      } else {
        wrapper.appendFormatted(trimmedProof,
            prefix: "[C]<font size='tall'>", suffix: "</font>");
      }
    }
    if (note.isNotEmpty) {
      wrapper.line(len: 2);
      wrapper.appendFormatted(note, prefix: "[C]", suffix: "");
    }
    return wrapper.print();
  }
}

class PurchaseBillModel {
  int id;
  int billAmount;
  int admin;
  int totalPrice;
  String customerName;
  Map<String, dynamic> data;
  List<Map<String, String>> billData;

  PurchaseBillModel(
      {required this.id,
      required this.billAmount,
      required this.admin,
      required this.totalPrice,
      required this.customerName,
      required this.data,
      required this.billData});

  factory PurchaseBillModel.fromJson(Map<String, dynamic> json) {
    return PurchaseBillModel(
        id: json["id"],
        billAmount: json["billAmount"],
        admin: json["admin"],
        totalPrice: json["totalPrice"],
        customerName: json["customerName"],
        data: json["data"],
        billData: (json["billData"] as List<dynamic>).map((data) {
          Map<String, String> ret = {};
          (data as Map<String, dynamic>).forEach((key, value) {
            ret[key] = value.toString();
          });
          return ret;
        }).toList());
  }
}
