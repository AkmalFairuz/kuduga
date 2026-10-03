import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/screens/balance/balance_deposit_screen.dart';
import 'package:finance_app/screens/balance/balance_transfer_screen.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../../screens/balance/balance_screen.dart';

class BalanceCard extends StatelessWidget {
  const BalanceCard({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
        height: 150.sp,
        decoration: BoxDecoration(
          color: ColorHelper.bg(context),
          boxShadow: [
            BoxShadow(
              color: ColorHelper.invertedBg(context).withOpacity(0.1),
              spreadRadius: 2,
              blurRadius: 5,
              offset: const Offset(0, 2),
            ),
          ],
          borderRadius: BorderRadius.circular(6.sp),
        ),
        padding: EdgeInsets.all(14.sp),
        child: _BalanceContent());
  }
}

class _BalanceContent extends StatefulWidget {
  @override
  _BalanceContentState createState() => _BalanceContentState();
}

class _BalanceContentState extends State<_BalanceContent> {
  bool isBalanceVisible = true;
  final AccountController accountController = AccountController.getInstance();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              "Saldo tersedia",
              style: TextStyle(
                fontSize: 16.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            SizedBox(width: 6.sp),
            GestureDetector(
              onTap: () {
                setState(() {
                  isBalanceVisible = !isBalanceVisible;
                });
              },
              child: Icon(
                isBalanceVisible ? Icons.visibility : Icons.visibility_off,
                size: 18.sp,
              ),
            )
          ],
        ),
        SizedBox(height: 3.sp),
        Expanded(
            child: GestureDetector(
                onTap: () {
                  Screens.to(const BalanceScreen());
                },
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Dyn(() {
                      // loading indicator if balance is null
                      if (!accountController.isFetched) {
                        return Skeleton(
                          width: 150.sp,
                          height: 30.sp,
                        );
                      }
                      // if balance is not null
                      if (accountController.balance != null) {
                        return Text(
                          isBalanceVisible
                              ? "Rp${Meta.currencyFormat(accountController.balance!)}"
                              : "Rp*****",
                          style: TextStyle(
                            fontSize: 24.sp,
                            fontWeight: FontWeight.w300,
                          ),
                        );
                      }
                      return null;
                    }, accountController)
                  ],
                ))),
        SizedBox(height: 10.sp),
        Row(
          children: [
            Expanded(
              child: Button(
                variant: ButtonVariant.outline,
                onPressed: () {
                  Screens.to(const BalanceTransferScreen());
                },
                child: const Text(
                  "Transfer",
                ),
              ),
            ),
            SizedBox(width: 8.sp),
            Expanded(
              child: Button(
                onPressed: () => Screens.to(const BalanceDepositScreen()),
                child: const Text(
                  "Deposit",
                ),
              ),
            ),
          ],
        )
      ],
    );
  }
}
