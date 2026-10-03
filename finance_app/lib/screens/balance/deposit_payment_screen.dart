import 'package:finance_app/model/deposit.dart';
import 'package:finance_app/screens/balance/deposit/deposit_webview_content.dart';
import 'package:finance_app/screens/balance/deposit_info_screen.dart';
import 'package:finance_app/screens/support/create_ticket_screen.dart';
import 'package:finance_app/service/balance_service.dart';
import 'package:finance_app/service/notification_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/utils/toast.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/fields/fields.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:finance_app/widgets/text/timer_text.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'deposit/bank_transfer_content.dart';
import 'deposit/qris_content.dart';

class DepositPaymentScreen extends StatefulWidget {
  const DepositPaymentScreen({Key? key, required this.deposit})
      : super(key: key);
  final DepositDetailedModel deposit;

  @override
  State<DepositPaymentScreen> createState() => _DepositPaymentState();
}

class _DepositPaymentState extends State<DepositPaymentScreen> {
  late DepositDetailedModel _deposit;

  @override
  void initState() {
    _deposit = widget.deposit;
    NotificationService.instance.subscribeMessage(handleNotification);

    super.initState();
  }

  @override
  void dispose() {
    NotificationService.instance.unsubscribeMessage(handleNotification);
    super.dispose();
  }

  void handleNotification(String name, Map<String, dynamic> data) {
    if (name != "deposit") {
      return;
    }
    int depositId = int.parse(data["depositId"]);
    if (depositId != _deposit.id) {
      return;
    }
    Alert.withLoading((done) async {
      final deposit = await BalanceService.getDepositDetailed(depositId);
      if (deposit.isSuccess()) {
        done();
        Screens.replace(DepositInfoScreen(deposit: deposit));
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_deposit.paymentUrl != null) {
      return DepositWebViewContent(url: _deposit.paymentUrl!);
    }
    return Scaffold(
        appBar: const TextAppBar("Detail Pembayaran"),
        backgroundColor: ColorHelper.theme(context,
            dark: Colors.grey[900], light: Colors.grey[200]!),
        body: RefreshIndicator(
          onRefresh: () async {
            DepositDetailedModel newDeposit =
                await BalanceService.getDepositDetailed(_deposit.id);
            if (newDeposit.isSuccess()) {
              Screens.replace(DepositInfoScreen(deposit: newDeposit));
              return;
            }
            setState(() {
              _deposit = newDeposit;
            });
          },
          child: _buildContentWrapper(context),
        ));
  }

  Widget _buildContentWrapper(BuildContext context) {
    final decor = BoxDecoration(
      color: ColorHelper.bg(context),
      borderRadius: BorderRadius.circular(8.sp),
      border: Border.all(
        color: ColorHelper.bg300(context),
        width: 1.sp,
      ),
    );

    return SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Column(
          children: [
            SizedBox(height: 8.sp),
            Container(
              margin: EdgeInsets.symmetric(horizontal: 16.sp, vertical: 8.sp),
              decoration: decor,
              child: Column(
                children: [
                  ...Fields.build(context, [
                    FieldEntry("ID Deposit", _deposit.id.toString()),
                    FieldEntry("Tanggal Kadaluarsa",
                        Meta.formatUnixDate(_deposit.expiredAt)),
                    FieldEntry(
                        "Kadaluarsa Dalam",
                        TimerText(
                            duration: Duration(
                                seconds: _deposit.expiredAt - Meta.unix()))),
                    FieldEntry("Metode Pembayaran", _deposit.paymentMethod),
                  ]),
                ],
              ),
            ),
            Container(
                padding:
                    EdgeInsets.symmetric(horizontal: 20.sp, vertical: 16.sp),
                decoration: decor,
                margin: EdgeInsets.symmetric(horizontal: 16.sp, vertical: 8.sp),
                child: _buildContent(context)),
            SizedBox(height: 8.sp),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.sp),
              child: Row(
                children: [
                  Expanded(
                      child: Button(
                          onPressed: _onCancel,
                          variant: ButtonVariant.outline,
                          child: const Text("Batalkan"))),
                  SizedBox(width: 8.sp),
                  Expanded(
                      child: Button(
                          onPressed: _onNeedHelp,
                          child: const Text("Butuh Bantuan")))
                ],
              ),
            ),
            SizedBox(height: 16.sp),
          ],
        ));
  }

  void _onNeedHelp() {
    Screens.to(
        CreateTicketScreen(initialMessage: "Deposit ID: ${_deposit.id}\n"));
  }

  void _onCancel() async {
    if ((await Alert.showYesOrNo(
            context: context,
            title: "Batalkan Deposit",
            message: "Apakah anda ingin membatalkan deposit ini?")) !=
        true) {
      return;
    }
    await Alert.withLoading((done) async {
      await BalanceService.cancelDeposit(_deposit.id);
      done();
      if (mounted) {
        Toast.success(context, "Deposit #${_deposit.id} telah dibatalkan");
        Screens.back();
      }
    });
  }

  Widget _buildContent(BuildContext context) {
    switch (_deposit.paymentMethod) {
      case "Bank BCA":
        return BankTransferContent(_deposit);
      case "QRIS":
        return QrisContent(_deposit);
    }
    return const Center(
      child: Text("Invalid payment method"),
    );
  }
}
