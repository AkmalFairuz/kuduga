import 'package:finance_app/service/local_auth_service.dart';
import 'package:finance_app/utils/bottom_select.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/button/list_button_view.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../themes.dart';
import '../../utils/alert.dart';

class AppSettingsScreen extends StatelessWidget {
  const AppSettingsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final systemThemeStr =
        MediaQuery.of(context).platformBrightness == Brightness.dark
            ? "Gelap"
            : "Terang";
    return Scaffold(
      appBar: const TextAppBar("Pengaturan Aplikasi"),
      body: ListButtonView(
        buttons: [
          ListButton(
              title: "Tema",
              onPress: () async {
                final result = await BottomSelect.show(
                    context: context,
                    selectedValue: Themes.current,
                    title: "Tema",
                    items: [
                      BottomSelectItem(
                          value: "system",
                          widget: Text("Sistem ($systemThemeStr)")),
                      BottomSelectItem(
                          value: "light", widget: const Text("Terang")),
                      BottomSelectItem(
                          value: "dark", widget: const Text("Gelap")),
                    ]);
                if (result != null) {
                  Themes.changeThemeMode(Themes.modeFromString(result));
                }
              }),
          ListButton(
              title: "Kunci Fingerprint",
              onPress: () async {
                final result = await BottomSelect.show(
                    context: context,
                    title: "Kunci Fingerprint",
                    selectedValue:
                        await LocalAuthService.isFingerprintLockEnabled(),
                    items: [
                      BottomSelectItem(
                          value: true, widget: const Text("Aktif")),
                      BottomSelectItem(
                          value: false, widget: const Text("Nonaktif")),
                    ]);
                if (result != null) {
                  try {
                    await LocalAuthService.setFingerprintLockEnabled(result);
                    if (result) {
                      Alert.message("Kunci fingerprint berhasil diaktifkan");
                    } else {
                      Alert.message("Kunci fingerprint berhasil dinonaktifkan");
                    }
                  } catch (e) {
                    if (e is Exception) {
                      Alert.error(e);
                    }
                  }
                }
              }),
        ],
      ),
    );
  }
}
