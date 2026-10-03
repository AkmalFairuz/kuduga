import 'package:finance_app/constants.dart';
import 'package:finance_app/screens/auth/agreement_screen.dart';
import 'package:finance_app/screens/auth/login_screen.dart';
import 'package:finance_app/screens/auth/register_screen.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:flutter/material.dart';

class IntroScreen extends StatelessWidget {
  const IntroScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          const SizedBox(height: 32),
          const Text("Selamat datang di Kuduga",
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
          const SizedBox(height: 16),
          const Text(
            "Kuduga menyediakan layanan pembelian pulsa, token listrik, voucher game dan produk digital lainnya.",
            textAlign: TextAlign.center,
          ),
          Expanded(child: Container()),
          const SizedBox(height: 16),
          Button(
              onPressed: () async {
                Screens.to(AgreementScreen(
                  onAgree: () => Screens.replace(const RegisterScreen()),
                ));
              },
              width: double.infinity,
              padding: const EdgeInsets.symmetric(
                vertical: 10,
              ),
              child: const Text("Daftar")),
          const SizedBox(height: 8),
          Button(
              variant: ButtonVariant.outline,
              onPressed: () {
                Screens.to(
                  const LoginScreen(),
                );
              },
              width: double.infinity,
              padding: const EdgeInsets.symmetric(
                vertical: 10,
              ),
              child: const Text("Login")),
          const SizedBox(height: 12),
          const Text("v${Constants.version}", style: TextStyle(fontSize: 10.5)),
        ],
      ),
    ));
  }
}
