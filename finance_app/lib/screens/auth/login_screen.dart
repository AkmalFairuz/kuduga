import 'package:finance_app/screens/auth/agreement_screen.dart';
import 'package:finance_app/screens/auth/register_screen.dart';
import 'package:finance_app/screens/auth/reset_password_screen.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/service/auth_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/text/text_link.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../widgets/button/button.dart';
import '../../widgets/input/input.dart';
import '../../widgets/text/text_app_bar.dart';
import '../home_screen.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Login"),
      body: SingleChildScrollView(
        child: Container(
          padding: const EdgeInsets.all(16),
          child: const _LoginForm(),
        ),
      ),
    );
  }
}

class _LoginForm extends StatefulWidget {
  const _LoginForm();

  @override
  _LoginFormState createState() => _LoginFormState();
}

class _LoginFormState extends State<_LoginForm> {
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  bool isPasswordVisible = false;

  void _handleLogin() {
    // validate fields
    if (_usernameController.text.isEmpty) {
      Alert.message("Username atau email tidak boleh kosong");
      return;
    }
    if (_passwordController.text.isEmpty) {
      Alert.message("Password tidak boleh kosong");
      return;
    }
    Alert.withLoading((done) async {
      final token = await AuthService.login(
          _usernameController.text, _passwordController.text);
      done();
      _handleLoginSuccess(token);
    });
  }

  void _handleLoginSuccess(String token) {
    //Toast.success(context, 'Login berhasil');
    Server.setToken(token);
    Screens.replaceAll(const HomeScreen());
  }

  void _handleLoginWithGoogle() async {
    try {
      final token = await AuthService.signInWithGoogle();
      _handleLoginSuccess(token);
    } catch (e) {
      Alert.message(e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Input(
              controller: _usernameController,
              hintText: "Masukkan Username atau Alamat Email",
              label: "Username atau Alamat Email",
              prefixIcon: const Icon(Icons.person),
            ),
            const SizedBox(height: 16),
            Input(
              controller: _passwordController,
              hintText: "Masukkan Password",
              label: "Password",
              isPassword: !isPasswordVisible,
              prefixIcon: const Icon(Icons.lock),
              suffixIcon: IconButton(
                icon: Icon(isPasswordVisible
                    ? Icons.visibility
                    : Icons.visibility_off),
                onPressed: () {
                  setState(() {
                    isPasswordVisible = !isPasswordVisible;
                  });
                },
              ),
            ),
            const SizedBox(height: 8),
            Align(
              alignment: Alignment.centerRight,
              child: TextLink(
                text: "Lupa Password?",
                onPressed: () => Screens.to(const ResetPasswordScreen()),
              ),
            ),
            const SizedBox(height: 8),
            Button(
              width: double.infinity,
              onPressed: _handleLogin,
              child: const Text("Login"),
            ),
            const SizedBox(height: 16),
            const Text("atau"),
            const SizedBox(height: 16),
            Button(
                width: double.infinity,
                variant: ButtonVariant.white,
                onPressed: _handleLoginWithGoogle,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    SvgPicture.asset("assets/icons/google.svg",
                        width: 20, height: 20),
                    const SizedBox(width: 10),
                    const Text("Login dengan Google"),
                  ],
                )),
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text("Belum punya akun?"),
                const SizedBox(width: 4),
                TextLink(
                    text: "Daftar",
                    onPressed: () {
                      Screens.replace(AgreementScreen(
                          onAgree: () =>
                              Screens.replace(const RegisterScreen())));
                    }),
              ],
            ),
          ],
        ),
      ],
    );
  }
}
