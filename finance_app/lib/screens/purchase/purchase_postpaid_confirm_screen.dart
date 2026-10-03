import 'package:finance_app/model/product.dart';
import 'package:finance_app/model/purchase.dart';
import 'package:finance_app/service/purchase_service.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../controller/account_controller.dart';
import '../auth/pin_screen.dart';

class PurchasePostpaidConfirmScreen extends StatelessWidget {
  const PurchasePostpaidConfirmScreen(
      {required this.bill,
      required this.product,
      required this.productDestination,
      required this.destination,
      Key? key})
      : super(key: key);

  final PurchaseBillModel bill;
  final ProductModel product;
  final ProductDestinationModel productDestination;
  final List<ProductDestinationFieldRequestModel> destination;

  List<Widget> _buildBillData2(
      BuildContext context, int index, Map<String, String> billData) {
    List<Widget> ret = [];
    ret.add(Container(
      color: Meta.color,
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
      child: Text('Detail Tagihan #${index + 1}',
          style: const TextStyle(
              color: Colors.white, fontWeight: FontWeight.bold)),
    ));
    ret.addAll(Fields.build(
        context,
        billData.entries
            .map((entry) => FieldEntry(entry.key, entry.value.toString()))
            .toList()));
    return ret;
  }

  List<Widget> _buildBillData(
    BuildContext context,
  ) {
    List<Widget> ret = [];
    int i = 0;
    for (var data in bill.billData) {
      ret.addAll(_buildBillData2(context, i, data));
      i++;
    }
    return ret;
  }

  Future<void> _handleConfirm(BuildContext context) async {
    String? pin;
    if (AccountController.getInstance().authDetails!.hasPin) {
      pin = await Screens.to(
          const PinScreen(title: "Masukkan PIN untuk membayar tagihan"));
      if (pin == null) {
        return;
      }
    }

    Alert.withLoading((done) async {
      final purchaseId =
          await PurchaseService.payBill(product.id, bill.id, destination, pin);
      done();
      StateNotifier.notify("purchase");
      if (context.mounted) {
        Navigator.of(context).pop(purchaseId);
      }
    }, message: "Membayar tagihan...");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Konfirmasi"),
      body: Column(
        children: [
          Expanded(
              child: ListView(
            children: [
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  product.hasImage()
                      ? Image.network(product.imageUrl, width: 40, height: 40)
                      : const Icon(Icons.question_mark, size: 40),
                  const SizedBox(width: 16),
                  Text(product.name)
                ],
              ),
              const SizedBox(height: 16),
              ...Fields.build(context, [
                FieldEntry(
                    "Jumlah Tagihan", Meta.currencyFormatRp(bill.billAmount)),
                FieldEntry("Biaya Admin", Meta.currencyFormatRp(bill.admin)),
                FieldEntry(
                    "Total Biaya", Meta.currencyFormatRp(bill.totalPrice)),
                FieldEntry("Nama", bill.customerName),
                ...bill.data.entries
                    .map((entry) =>
                        FieldEntry(entry.key, entry.value.toString()))
                    .toList(),
              ]),
              ..._buildBillData(context)
            ],
          )),
          const Line(),
          Container(
            padding: const EdgeInsets.all(16),
            child: Button(
                onPressed: () => _handleConfirm(context),
                width: double.infinity,
                child: const Text('Bayar Sekarang')),
          )
        ],
      ),
    );
  }
}
