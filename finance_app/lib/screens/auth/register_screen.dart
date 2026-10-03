import 'package:finance_app/model/register.dart';
import 'package:finance_app/screens/auth/login_screen.dart';
import 'package:finance_app/screens/auth/register_otp_screen.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/service/auth_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';

import '../../widgets/text/text_app_bar.dart';
import '../../widgets/text/text_link.dart';

class RegisterScreen extends StatelessWidget {
  const RegisterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Daftar"),
        body: SingleChildScrollView(
          child: Container(
            padding: const EdgeInsets.all(16),
            child: const _RegisterForm(),
          ),
        ));
  }
}

class _RegisterForm extends StatefulWidget {
  const _RegisterForm();

  @override
  _RegisterFormState createState() => _RegisterFormState();
}

class _RegisterFormState extends State<_RegisterForm> {
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  bool isPasswordVisible = false;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Input(
                prefixIcon: const Icon(Icons.person),
                label: 'Nama Lengkap',
                controller: _nameController,
                hintText: 'Masukkan nama lengkap'),
            const SizedBox(height: 10),
            Input(
              prefixIcon: const Icon(Icons.alternate_email),
              label: 'Username',
              controller: _usernameController,
              hintText: 'Masukkan username',
            ),
            const SizedBox(height: 10),
            Input(
              prefixIcon: const Icon(Icons.email),
              label: 'Alamat Email',
              controller: _emailController,
              hintText: 'Masukkan alamat email',
            ),
            const SizedBox(height: 10),
            Input(
              prefixIcon: const Icon(Icons.lock),
              controller: _passwordController,
              hintText: 'Masukkan kata sandi',
              label: 'Kata Sandi',
              isPassword: !isPasswordVisible,
              suffixIcon: IconButton(
                splashRadius: 1,
                icon: Icon(isPasswordVisible
                    ? Icons.visibility_off
                    : Icons.visibility),
                onPressed: () {
                  setState(() {
                    isPasswordVisible = !isPasswordVisible;
                  });
                },
              ),
            ),
          ],
        ),
        const SizedBox(height: 20),
        Button(
            onPressed: _handleFormSubmit,
            width: double.infinity,
            child: const Text('Daftar')),
        const SizedBox(height: 16),
        const Center(
          child: Text(
            'atau',
          ),
        ),
        const SizedBox(height: 16),
        Button(
            width: double.infinity,
            variant: ButtonVariant.white,
            onPressed: _registerWithGoogle,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SvgPicture.asset("assets/icons/google.svg",
                    width: 20, height: 20),
                const SizedBox(width: 10),
                const Text("Daftar dengan Google"),
              ],
            )),
        const SizedBox(height: 32),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text("Sudah mempunyai akun?"),
            const SizedBox(width: 4),
            TextLink(
                text: 'Login',
                onPressed: () => Screens.replace(const LoginScreen())),
          ],
        ),
      ],
    );
  }

  void _registerWithGoogle() async {
    final token = await AuthService.signUpWithGoogle();
    await Server.setToken(token);
    Meta.afterRegister();
  }

  void _handleFormSubmit() async {
    RegisterModel model = RegisterModel(
        displayName: _nameController.text,
        username: _usernameController.text.toLowerCase(),
        email: _emailController.text,
        password: _passwordController.text);

    try {
      model.validate();
    } catch (e) {
      Alert.message(e.toString());
      return;
    }

    Alert.withLoading((done) async {
      final otp = await AuthService.getRegisterOTP(model);
      model.otp = otp;
      done();
      Screens.to(RegisterOTPScreen(model: model));
    });
  }
}
