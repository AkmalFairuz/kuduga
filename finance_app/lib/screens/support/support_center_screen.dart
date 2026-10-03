import 'package:finance_app/screens/support/create_ticket_screen.dart';
import 'package:finance_app/screens/support/ticket_list_screen.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/button/list_button_view.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class SupportCenterScreen extends StatelessWidget {
  const SupportCenterScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Pusat Bantuan"),
      body: ListButtonView(
        buttons: [
          ListButton(
              title: "Lihat Tiket Bantuan Anda",
              onPress: () => Screens.to(const TicketListScreen())),
          ListButton(
              title: "Buat Tiket Bantuan",
              onPress: () {
                Screens.to(const CreateTicketScreen());
              }),
        ],
      ),
    );
  }
}
