import 'package:finance_app/model/transfer.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class BalanceTransferInfoScreen extends StatelessWidget {
  const BalanceTransferInfoScreen({Key? key, required this.model})
      : super(key: key);

  final TransferModel model;

  @override
  Widget build(BuildContext context) {
    var appBar = const TextAppBar("Detail Transfer");
    return Scaffold(
        appBar: appBar,
        body: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const SizedBox(height: 32),
              Container(
                width: 110,
                height: 110,
                decoration: BoxDecoration(
                    color: Meta.color,
                    borderRadius: BorderRadius.circular(999)),
                child: Center(
                  child: Icon(
                    Icons.check,
                    size: 80,
                    color: Colors.grey[50],
                  ),
                ),
              ),
              const SizedBox(height: 32),
              const Text(
                "Transfer berhasil",
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 12),
              Text(model.isSend
                  ? "Anda telah melakukan transfer sebesar:"
                  : "Anda menerima transfer sebesar:"),
              const SizedBox(height: 12),
              Text(
                Meta.currencyFormatRp(model.amount),
                style:
                    const TextStyle(fontWeight: FontWeight.bold, fontSize: 24),
              ),
              const SizedBox(height: 32),
              ...Fields.build(context, [
                FieldEntry("ID", "${model.id}"),
                FieldEntry(
                    "Pengirim", "${model.sender.name}\n${model.sender.email}"),
                FieldEntry("Penerima",
                    "${model.receiver.name}\n${model.receiver.email}"),
                FieldEntry("Nominal", Meta.currencyFormatRp(model.amount)),
                FieldEntry("Catatan", model.note),
                FieldEntry("Tanggal", Meta.formatUnixDate(model.createdAt)),
              ]),
            ],
          ),
        ));
  }
}
