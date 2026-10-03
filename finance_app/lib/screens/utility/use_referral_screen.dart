import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/local_storage.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class UseReferralScreen extends StatefulWidget {
  const UseReferralScreen({super.key});

  @override
  State<UseReferralScreen> createState() => _UseReferralScreenState();
}

class _UseReferralScreenState extends State<UseReferralScreen> {
  final referralController = TextEditingController();

  void onContinue() async {
    if (referralController.text.isNotEmpty) {
      await Alert.withLoading((done) async {
        final message = await UserService.useReferral(referralController.text);
        done();
        Screens.back();
        LocalStorage.set("alreadyUsedReferral", true);
        await Alert.message(message);
      });
      return;
    }
    Screens.back();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Gunakan Kode Referral"),
        body: Container(
          padding: const EdgeInsets.all(16),
          height: double.infinity,
          width: double.infinity,
          child: Column(
            children: [
              const Text(
                "Jika Anda memiliki kode referral, silakan masukkan kode referral Anda disini. Kosongkan jika Anda tidak memiliki kode referral.",
              ),
              const SizedBox(height: 16),
              Input(
                controller: referralController,
                hintText: "Kode referral (opsional)",
              ),
              const SizedBox(height: 16),
              Button(
                  onPressed: onContinue,
                  width: double.infinity,
                  child: const Text("Lanjut")),
            ],
          ),
        ));
  }
}
