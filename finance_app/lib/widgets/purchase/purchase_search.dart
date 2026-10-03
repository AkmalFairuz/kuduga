import 'package:finance_app/screens/purchase/purchase_search_result_screen.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/styles/text_style_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

class PurchaseSearch extends StatefulWidget {
  const PurchaseSearch({Key? key}) : super(key: key);

  @override
  State<PurchaseSearch> createState() => _PurchaseSearchState();
}

class _PurchaseSearchState extends State<PurchaseSearch> {
  final queryController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      const SizedBox(height: 16),
      Input(
        controller: queryController,
        label: 'Tujuan',
      ),
      const SizedBox(height: 4),
      Text(
          "Tujuan bisa berisi Nomor HP, ID Pelanggan PLN, Nomor Internet Indihome dan lainnya",
          style: TextStyleHelper.inputHelper(context)),
      const SizedBox(height: 16),
      Button(
          onPressed: () async {
            if (queryController.text.length < 3) {
              await Alert.message("Minimal 3 karakter");
              return;
            }
            Alert.withLoading((done) async {
              final result =
                  await PurchaseService.searchPurchases(queryController.text);
              done();
              if (result.isEmpty) {
                await Alert.message(
                    "Tidak ada hasil dari pencarian '${queryController.text}'");
                return;
              }
              await Screens.to(PurchaseSearchResultScreen(
                  query: queryController.text, result: result));
              Screens.back();
            });
          },
          width: double.infinity,
          child: const Text("Cari"))
    ]);
  }
}
