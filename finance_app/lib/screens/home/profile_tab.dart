import 'package:finance_app/constants.dart';
import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/screens/balance/balance_screen.dart';
import 'package:finance_app/screens/settings/account_security_screen.dart';
import 'package:finance_app/screens/settings/app_settings_screen.dart';
import 'package:finance_app/screens/settings/personal_settings_screen.dart';
import 'package:finance_app/screens/stats/stats_screen.dart';
import 'package:finance_app/screens/support/support_center_screen.dart';
import 'package:finance_app/screens/utility/about_us_screen.dart';
import 'package:finance_app/screens/utility/referral_info_screen.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class ProfileTab extends StatelessWidget {
  const ProfileTab({Key? key, required this.onOpenPurchaseHistory})
      : super(key: key);

  final VoidCallback onOpenPurchaseHistory;

  @override
  Widget build(BuildContext context) {
    Widget separator = Container(height: 20, color: ColorHelper.bg100(context));
    AccountController acc = AccountController.getInstance();

    List<Widget> buttons = [
      Container(
          color: Theme.of(context).appBarTheme.backgroundColor,
          padding: const EdgeInsets.all(20),
          child: DefaultTextStyle(
              style: Theme.of(context).textTheme.bodyLarge!.copyWith(
                    color: Colors.white,
                  ),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.white, width: 2),
                      gradient: LinearGradient(colors: [
                        Meta.color[300]!,
                        Meta.color[900]!,
                      ]),
                      borderRadius: BorderRadius.circular(999),
                    ),
                    child: Icon(
                      Icons.person,
                      size: 48,
                      color: Colors.white,
                    ),
                  ),
                  const SizedBox(width: 20),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Dyn(
                          () => Text(
                              acc.authDetails != null
                                  ? acc.authDetails!.displayName
                                  : '...',
                              style: const TextStyle(
                                  fontWeight: FontWeight.bold, fontSize: 20)),
                          acc),
                      Dyn(
                          () => Text(
                              '@${acc.authDetails != null ? acc.authDetails!.username : ''}'),
                          acc),
                      const SizedBox(height: 4),
                      Row(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          Container(
                              padding: const EdgeInsets.symmetric(
                                  vertical: 2, horizontal: 6),
                              decoration: BoxDecoration(
                                  color: Meta.color[600]!,
                                  borderRadius: BorderRadius.circular(4)),
                              child: Dyn(
                                  () => Text(
                                      acc.balance != null
                                          ? Meta.currencyFormatRp(acc.balance!)
                                          : 'Loading...',
                                      style: const TextStyle(
                                          fontSize: 11,
                                          fontWeight: FontWeight.w600)),
                                  acc)),
                        ],
                      )
                    ],
                  )
                ],
              ))),
      ListButton(
          icon: Icons.account_balance_wallet_outlined,
          title: 'Saldo',
          onPress: () => Screens.to(const BalanceScreen())),
      ListButton(
          icon: Icons.history,
          title: 'Riwayat Pembelian',
          onPress: () {
            onOpenPurchaseHistory();
          }),
      ListButton(
        icon: Icons.show_chart,
        title: 'Statistik',
        onPress: () {
          Screens.to(const StatsScreen());
        },
      ),
      ListButton(
          icon: CupertinoIcons.person,
          title: 'Pengaturan Pribadi',
          onPress: () {
            Screens.to(const PersonalSettingsScreen());
          }),
      ListButton(
          icon: CupertinoIcons.person_add,
          title: "Undang Teman",
          onPress: () {
            Screens.to(ReferralInfoScreen());
          }),
      ListButton(
          icon: CupertinoIcons.lock,
          title: 'Keamanan Akun',
          onPress: () {
            Screens.to(const AccountSecurityScreen());
          }),
      separator,
      ListButton(
          icon: CupertinoIcons.question_circle,
          title: 'Pusat Bantuan',
          onPress: () {
            Screens.to(const SupportCenterScreen());
          }),
      ListButton(
          icon: CupertinoIcons.info,
          title: 'Tentang Kami',
          onPress: () async {
            Screens.to(const AboutUsScreen());
          }),
      separator,
      // app settings
      ListButton(
          icon: Icons.settings,
          title: 'Pengaturan Aplikasi',
          onPress: () {
            Screens.to(const AppSettingsScreen());
          }),
      ListButton(
          icon: Icons.logout,
          title: 'Logout',
          onPress: () {
            Alert.showYesOrNo(
              context: context,
              title: 'Logout',
              message: 'Apakah anda yakin ingin logout?',
              onYes: () {
                Server.logout(context);
              },
            );
          }),
      Container(
        color: ColorHelper.bg100(context),
        width: double.infinity,
        padding: const EdgeInsets.all(20),
        child: const Column(
          children: [
            Text(Constants.version),
            SizedBox(height: 4),
          ],
        ),
      ),
    ];

    return Container(
      color: ColorHelper.bg100(context),
      height: double.infinity,
      child: SingleChildScrollView(
        child: Container(
            color: Theme.of(context).scaffoldBackgroundColor,
            child: Column(
              children: Meta.addSeparatorToList(buttons, const Line()),
            )),
      ),
    );
  }
}
