import 'package:finance_app/screens/utility/webview_screen.dart';
import 'package:finance_app/service/app_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/button/list_button_view.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class AboutUsScreen extends StatelessWidget {
  const AboutUsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Tentang Kami"),
      body: ListButtonView(
        buttons: [
          ListButton(
              title: "Tentang Kami",
              onPress: () async {
                Screens.to(const WebviewScreen("https://example.invalid/about"));
              }),
          ListButton(
              title: "Syarat dan Ketentuan",
              onPress: () async {
                Alert.withLoading((done) async {
                  final html = await AppService.getTermsAndConditions();
                  done();
                  Screens.to(WebviewScreen(html, type: WebviewScreenType.html));
                });
              }),
          ListButton(
              title: "Kebijakan Privasi",
              onPress: () async {
                Alert.withLoading((done) async {
                  final html = await AppService.getPrivacyPolicy();
                  done();
                  Screens.to(WebviewScreen(html, type: WebviewScreenType.html));
                });
              }),
        ],
      ),
    );
  }
}
