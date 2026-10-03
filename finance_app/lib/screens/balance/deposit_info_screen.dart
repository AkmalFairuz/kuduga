import 'package:finance_app/model/deposit.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class DepositInfoScreen extends StatelessWidget {
  const DepositInfoScreen({Key? key, required this.deposit}) : super(key: key);

  final DepositDetailedModel deposit;

  @override
  Widget build(BuildContext context) {
    Color? statusColor;
    switch (deposit.status) {
      case DepositStatus.waiting:
        statusColor = Colors.orange;
        break;
      case DepositStatus.success:
        statusColor = Colors.green[700]!;
        break;
      case DepositStatus.failed:
        statusColor = Colors.red;
        break;
    }

    IconData? statusIcon;
    switch (deposit.status) {
      case DepositStatus.waiting:
        statusIcon = Icons.pending;
        break;
      case DepositStatus.success:
        statusIcon = Icons.check;
        break;
      case DepositStatus.failed:
        statusIcon = Icons.close;
        break;
    }

    return Scaffold(
      appBar: TextAppBar("Deposit #${deposit.id}"),
      body: SingleChildScrollView(
        child: Column(
          children: [
            Stack(children: [
              Column(
                children: [
                  const SizedBox(height: 100),
                  ...Fields.build(context, [
                    FieldEntry("ID Deposit", deposit.id.toString()),
                    FieldEntry("Metode Pembayaran", deposit.paymentMethod),
                    FieldEntry(
                        "Nominal", Meta.currencyFormatRp(deposit.amount)),
                    FieldEntry(
                        "Biaya Admin", Meta.currencyFormatRp(deposit.fee)),
                    FieldEntry(
                        "Status", DepositStatus.displayStatus(deposit.status)),
                    FieldEntry(
                        "Tanggal", Meta.formatUnixDate(deposit.createdAt)),
                  ])
                ],
              ),
              // icon status above container in center
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    margin: const EdgeInsets.only(top: 16),
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                        color: statusColor,
                        borderRadius: BorderRadius.circular(999)),
                    child: Icon(statusIcon, size: 40, color: Colors.white),
                  )
                ],
              ),
            ]),
            const SizedBox(height: 16),
            const Line(),
            _buildTrackWidgets(),
          ],
        ),
      ),
    );
  }

  Widget _buildTrackWidgets() {
    return Column(
      children:
          deposit.tracks.map((track) => _buildTrackWidget(track)).toList(),
    );
  }

  Widget _buildTrackWidget(DepositTrack track) {
    return Container(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Flexible(
                  child: Text(track.description),
                ),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(Meta.formatUnixDate(track.createdAt),
                        style: const TextStyle(fontSize: 12)),
                    const SizedBox(height: 2),
                    if (track.newStatus != null)
                      Text(DepositStatus.displayStatus(track.newStatus!),
                          style: const TextStyle(
                              fontSize: 11, fontWeight: FontWeight.bold)),
                  ],
                )
              ],
            )
          ],
        ));
  }
}
