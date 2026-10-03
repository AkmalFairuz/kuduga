import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/screens/settings/change_display_name_screen.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/material.dart';

import '../../widgets/button/button.dart';
import '../../widgets/text/text_app_bar.dart';
import '../kyc/kyc_screen.dart';
import 'change_email_screen.dart';

class PersonalSettingsScreen extends StatelessWidget {
  const PersonalSettingsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    AccountController controller = AccountController.getInstance();
    return Scaffold(
        appBar: const TextAppBar("Pengaturan Pribadi"),
        body: SingleChildScrollView(
            child: Container(
          padding: const EdgeInsets.all(16),
          child: Dyn(
            () => Column(
              children: [
                Input(
                    enabled: false,
                    label: 'Username',
                    controller: TextEditingController(
                        text: controller.authDetails!.username)),
                const SizedBox(height: 16),
                Input(
                    enabled: false,
                    suffixIcon: GestureDetector(
                        onTap: () {
                          Screens.to(const ChangeEmailScreen());
                        },
                        child: const Icon(Icons.edit)),
                    label: 'Alamat Email',
                    controller: TextEditingController(
                        text: controller.authDetails!.email)),
                const SizedBox(height: 16),
                Input(
                    enabled: false,
                    suffixIcon: GestureDetector(
                        onTap: () {
                          Screens.to(const ChangeDisplayNameScreen());
                        },
                        child: const Icon(Icons.edit)),
                    label: 'Nama Lengkap',
                    controller: TextEditingController(
                        text: controller.authDetails!.displayName)),
                const SizedBox(height: 16),
                FutureBuilder(
                    future: UserService.getKycStatus(),
                    builder: (context, snapshot) {
                      if (!snapshot.hasData) {
                        return const Skeleton(
                          width: double.infinity,
                          height: 50,
                        );
                      }
                      final data = snapshot.data!;
                      if (data.currentStatus) {
                        return Row(
                          children: [
                            Icon(
                              Icons.check,
                              size: 24,
                              color: Colors.green[600],
                            ),
                            const SizedBox(width: 16),
                            Text(
                              "Identitas terverifikasi",
                              style: TextStyle(
                                  color: Colors.green[600],
                                  fontSize: 15,
                                  fontWeight: FontWeight.bold),
                            ),
                          ],
                        );
                      }
                      if (data.requestStatus == null) {
                        return Button(
                          onPressed: () {
                            Screens.to(const KycScreen());
                          },
                          width: double.infinity,
                          child: const Text("Verifikasi Identitas"),
                        );
                      }
                      return const Row(
                        children: [
                          Icon(
                            Icons.access_time,
                            size: 24,
                          ),
                          SizedBox(width: 16),
                          Text(
                            "Identitas sedang dalam proses verifikasi",
                            style: TextStyle(
                                fontSize: 15, fontWeight: FontWeight.bold),
                          ),
                        ],
                      );
                    })
              ],
            ),
            controller,
          ),
        )));
  }
}
