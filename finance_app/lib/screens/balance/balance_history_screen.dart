import 'package:finance_app/model/transaction.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class BalanceHistoryScreen extends StatelessWidget {
  const BalanceHistoryScreen(
      {Key? key,
      required this.transactions,
      required this.startDate,
      required this.endDate,
      required this.historyDate})
      : super(key: key);

  final List<TransactionModel> transactions;
  final DateTime startDate;
  final DateTime endDate;
  final DateTime historyDate;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const AppBarTitle("Mutasi Saldo"),
      ),
      body: SafeArea(
          child: Column(
        children: [
          transactions.isEmpty
              ? Container(
                  padding: const EdgeInsets.all(16),
                  child: const Center(
                    child: Text("Tidak ada transaksi"),
                  ),
                )
              : Expanded(
                  child: ListView.separated(
                      itemBuilder: (context, index) {
                        TransactionModel trx = transactions[index];
                        return _BalanceStatement(trx);
                      },
                      separatorBuilder: (_, __) => const Line(),
                      itemCount: transactions.length),
                )
        ],
      )),
    );
  }
}

class _BalanceStatement extends StatelessWidget {
  const _BalanceStatement(this.model, {Key? key}) : super(key: key);

  final TransactionModel model;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Row(children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(model.description),
              const SizedBox(height: 4),
              Text(
                  "${Meta.formatUnixDate(model.createdAt, format: "dd MMM yyyy HH:mm")}\ndari ${Meta.currencyFormatRp(model.beforeBalance)} ke ${Meta.currencyFormatRp(model.afterBalance)}",
                  style: TextStyle(
                    color: Colors.grey[500],
                    fontSize: 11,
                  )),
            ],
          ),
        ),
        const SizedBox(width: 12),
        Text(
          Meta.currencyFormatRp(model.amount),
          style: TextStyle(
            color: model.amount >= 0 ? Colors.green : Colors.red,
            fontWeight: FontWeight.bold,
            fontSize: 14,
          ),
        ),
      ]),
    );
  }
}
