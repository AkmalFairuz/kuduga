import 'package:finance_app/model/purchase.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/purchase/purchase_list.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class PurchaseSearchResultScreen extends StatelessWidget {
  const PurchaseSearchResultScreen(
      {Key? key, required this.query, required this.result})
      : super(key: key);

  final String query;
  final List<PurchaseModel> result;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Hasil Pencarian"),
      body: Column(
        children: [
          Container(
              padding: const EdgeInsets.all(16),
              child:
                  Text("Menampilkan ${result.length} hasil untuk \"$query\"")),
          const Line(),
          Expanded(child: PurchaseList.build(context, true, result))
        ],
      ),
    );
  }
}
