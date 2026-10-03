import 'package:finance_app/screens/settings/change_password_screen.dart';
import 'package:finance_app/screens/settings/pin_settings_screen.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

import '../../widgets/button/list_button.dart';
import '../../widgets/button/list_button_view.dart';
import '../../widgets/text/text_app_bar.dart';

class AccountSecurityScreen extends StatelessWidget {
  const AccountSecurityScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Keamanan Akun"),
      body: ListButtonView(
        buttons: [
          ListButton(
            title: 'Perbarui Kata Sandi',
            onPress: () {
              Screens.to(const ChangePasswordScreen());
            },
          ),
          ListButton(
            title: 'Pengaturan PIN',
            onPress: () async {
              Screens.to(PinSettingsScreen());
            },
          ),
        ],
      ),
    );
  }
}
