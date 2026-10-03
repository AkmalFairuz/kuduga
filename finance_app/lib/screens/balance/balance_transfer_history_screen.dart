import 'package:finance_app/screens/balance/balance_transfer_info_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../model/transfer.dart';

class BalanceTransferHistoryScreen extends StatefulWidget {
  const BalanceTransferHistoryScreen({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _BalanceTransferHistoryState();
}

class _BalanceTransferHistoryState extends State<BalanceTransferHistoryScreen> {
  List<TransferModel>? _transfers;

  @override
  void initState() {
    super.initState();

    BalanceService.getTransfers().then((v) {
      setState(() {
        _transfers = v;
      });
    });
  }

  Widget _buildContent() {
    if (_transfers == null) {
      return ListView();
    }
    if (_transfers!.isEmpty) {
      return const Center(
        child: Text("Tidak ada transfer"),
      );
    }
    return ListView.separated(
        itemBuilder: (context, index) {
          var tf = _transfers![index];
          return ListTile(
            onTap: () {
              Screens.to(BalanceTransferInfoScreen(model: tf));
            },
            title: Text(tf.isSend ? tf.receiver.name : tf.sender.name),
            subtitle: Text(
              Meta.formatUnixDate(tf.createdAt),
              style: const TextStyle(fontSize: 12),
            ),
            trailing: Text(
              (tf.isSend ? "-" : "+") + Meta.currencyFormatRp(tf.amount),
              style: TextStyle(
                  color: tf.isSend
                      ? ColorHelper.textSm(context)
                      : Colors.green[600],
                  fontWeight: FontWeight.bold,
                  fontSize: 16),
            ),
          );
        },
        separatorBuilder: (_, __) => const Line(),
        itemCount: _transfers!.length);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Riwayat Transfer"),
      body: _buildContent(),
    );
  }
}
