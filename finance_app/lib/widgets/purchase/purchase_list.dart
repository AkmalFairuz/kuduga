import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/purchase/purchase_tile.dart';
import 'package:flutter/material.dart';

import '../../model/purchase.dart';

class PurchaseList {
  static Widget build(
      BuildContext context, bool isLast, List<PurchaseModel> purchases) {
    Map<DateTime, List<PurchaseModel>> separated =
        _separatePurchasesByDate(purchases);

    if (purchases.isEmpty) {
      return const Center(
        child: Text("Tidak ada pembelian"),
      );
    }

    List<Widget> widgets = [];

    separated.forEach((dt, value) {
      widgets.add(Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        color: ColorHelper.bg100(context),
        child: Text(
          Meta.formatDateTime(dt, format: "EEEE, d MMM yyyy").toUpperCase(),
          style: TextStyle(
              fontWeight: FontWeight.bold, color: ColorHelper.textSm(context)),
        ),
      ));
      for (var purchase in value) {
        widgets.add(
            PurchaseTile(key: Key(purchase.purchaseId.toString()), purchase));
      }
    });

    if (!isLast) {
      widgets.add(const Padding(
        padding: EdgeInsets.all(32),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(
              width: 36,
              height: 36,
              child: CircularProgressIndicator(),
            )
          ],
        ),
      ));
    }

    return ListView.separated(
        physics: const AlwaysScrollableScrollPhysics(),
        itemBuilder: (_, index) => widgets[index],
        separatorBuilder: (_, __) => const Line(),
        itemCount: widgets.length);
  }

  static Map<DateTime, List<PurchaseModel>> _separatePurchasesByDate(
      List<PurchaseModel> models) {
    Map<DateTime, List<PurchaseModel>> separatedMap = {};

    for (PurchaseModel model in models) {
      DateTime date1 =
          DateTime.fromMillisecondsSinceEpoch(model.createdAt * 1000);
      DateTime date2 = DateTime(date1.year, date1.month, date1.day);

      if (!separatedMap.containsKey(date2)) {
        separatedMap[date2] = [];
      }

      separatedMap[date2]!.add(model);
    }

    return separatedMap;
  }
}
