import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/model/transaction.dart';
import 'package:finance_app/screens/balance/balance_deposit_screen.dart';
import 'package:finance_app/screens/balance/balance_history_screen.dart';
import 'package:finance_app/screens/balance/balance_transfer_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/button/list_button_view.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class BalanceScreen extends StatelessWidget {
  const BalanceScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Saldo"),
      body: SafeArea(
        child: SingleChildScrollView(child: _buildContent(context)),
      ),
    );
  }

  Widget _buildContent(BuildContext context) {
    return Column(
      children: [
        Container(
            padding: const EdgeInsets.all(16), child: const _BalanceHeader()),
        const SizedBox(height: 24),
        const Line(),
        const _BalanceContent(),
        const Line(),
      ],
    );
  }
}

class _BalanceHistoryForm extends StatefulWidget {
  const _BalanceHistoryForm({Key? key}) : super(key: key);

  @override
  State<_BalanceHistoryForm> createState() => _BalanceHistoryFormState();
}

class _BalanceHistoryFormState extends State<_BalanceHistoryForm> {
  DateTime fromDate = DateTime.now().add(const Duration(days: -7));
  DateTime toDate = DateTime.now();

  String formatDate(DateTime dt) {
    String day = dt.day.toString().padLeft(2, '0');
    String month = dt.month.toString().padLeft(2, '0');
    String year = dt.year.toString();
    return "$day/$month/$year";
  }

  void _showBalanceHistory() async {
    Alert.showLoading(context: context);

    List<TransactionModel> transactions;
    try {
      transactions = await BalanceService.getTransactions(
          from: (fromDate.toUtc().millisecondsSinceEpoch / 1000).round(),
          to: (toDate.toUtc().millisecondsSinceEpoch / 1000).round());
    } catch (e) {
      Alert.closeLoading();
      if (mounted) {
        Alert.show(
            title: "Gagal menampilkan mutasi saldo", message: e.toString());
      }
      return;
    }
    if (mounted) {
      Alert.closeLoading(context: context);
    }

    if (mounted) {
      Navigator.of(context).pop(); // close modal bottom sheet
    }

    Screens.to(BalanceHistoryScreen(
        transactions: transactions,
        startDate: fromDate,
        endDate: toDate,
        historyDate: DateTime.now()));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Mutasi Saldo',
          style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 16),
        GestureDetector(
          child: Input(
            enabled: false,
            prefixIcon: const Icon(Icons.calendar_today),
            label: "Dari tanggal",
            controller: TextEditingController(text: formatDate(fromDate)),
            onTap: () {
              showDatePicker(
                      helpText: "Dari tanggal",
                      context: context,
                      initialDate: fromDate,
                      firstDate: DateTime.now().add(const Duration(days: -90)),
                      lastDate: DateTime.now().add(const Duration(days: -1)))
                  .then((dt) {
                if (dt == null) return;
                setState(() {
                  fromDate = dt;
                });
              });
            },
          ),
        ),
        const SizedBox(height: 16),
        Input(
          enabled: false,
          prefixIcon: const Icon(Icons.calendar_today),
          label: "Sampai tanggal",
          controller: TextEditingController(text: formatDate(toDate)),
          onTap: () {
            showDatePicker(
                    helpText: "Sampai tanggal",
                    context: context,
                    initialDate: toDate,
                    firstDate: DateTime.now().add(const Duration(days: -89)),
                    lastDate: DateTime.now())
                .then((dt) {
              if (dt == null) return;
              setState(() {
                toDate = dt;
              });
            });
          },
        ),
        const SizedBox(height: 16),
        const Text(
            'Anda hanya bisa melihat mutasi 90 hari terakhir\n\nBatas mutasi yang muncul adalah 300 transaksi terakhir dari tanggal yang dipilih',
            style: TextStyle(color: Colors.grey)),
        const SizedBox(height: 16),
        Button(
            onPressed: _showBalanceHistory,
            width: double.infinity,
            child: const Text('Tampilkan')),
        const SizedBox(height: 16),
      ],
    );
  }
}

class _BalanceContent extends StatelessWidget {
  const _BalanceContent({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    Widget separator = Container(
      height: 1,
      color: Colors.grey[300],
    );

    return ListButtonView(
      buttons: [
        ListButton(
            title: 'Isi Saldo',
            icon: CupertinoIcons.plus,
            onPress: () {
              Screens.to(const BalanceDepositScreen());
            }),
        ListButton(
            icon: CupertinoIcons.arrow_right,
            title: 'Transfer Saldo',
            onPress: () {
              Screens.to(const BalanceTransferScreen());
            }),
        // separator
        ListButton(
            title: 'Mutasi Saldo',
            icon: Icons.history,
            onPress: () {
              ModalBottomSheet.show(
                  context,
                  (context) =>
                      const IntrinsicHeight(child: _BalanceHistoryForm()));
            }),
      ],
    );
  }
}

class _BalanceHeader extends StatelessWidget {
  const _BalanceHeader({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: ColorHelper.bg0(context),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: ColorHelper.bg300(context)),
        boxShadow: [
          BoxShadow(
            color: ColorHelper.bg400(context).withOpacity(0.2),
            blurRadius: 2,
            offset: const Offset(2, 2),
          )
        ],
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "Saldo tersedia",
                  style: TextStyle(
                    fontSize: 16,
                    color: ColorHelper.textSm(context),
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 4),
                Dyn(
                    () => Text(
                          Meta.currencyFormatRp(
                              AccountController.getInstance().balance!),
                          style: const TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                    AccountController.getInstance())
              ],
            ),
          ),
        ],
      ),
    );
  }
}
