import 'package:finance_app/model/purchase.dart';
import 'package:finance_app/screens/purchase/purchase_detail_screen.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

class PurchaseTile extends StatelessWidget {
  const PurchaseTile(this.model, {Key? key}) : super(key: key);

  final PurchaseModel model;

  @override
  Widget build(BuildContext context) {
    String dst = "";
    if (model.destination.length == 1) {
      dst += " - ${model.destination.entries.first.value}";
    }
    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onTap: () {
        Screens.to(PurchaseDetailScreen(purchaseId: model.purchaseId));
      },
      child: Container(
        padding: const EdgeInsets.all(16),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Container(
              width: 35,
              height: 35,
              decoration: BoxDecoration(
                color: model.statusColor(),
                borderRadius: BorderRadius.circular(99),
              ),
              child: Center(
                child: Icon(model.statusIcon(), color: Colors.grey[50]),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    model.productName,
                    style: TextStyle(
                        fontSize: 14,
                        color: ColorHelper.textLg(context),
                        fontWeight: FontWeight.bold),
                  ),
                  Text(
                    model.productCategoryName + dst,
                    style: TextStyle(
                        color: ColorHelper.textSm(context), fontSize: 11.5),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 16),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  Meta.currencyFormatRp(model.price),
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
