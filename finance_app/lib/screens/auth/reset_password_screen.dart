import 'package:finance_app/model/otp.dart';
import 'package:finance_app/model/reset_password.dart';
import 'package:finance_app/screens/auth/base_otp.dart';
import 'package:finance_app/screens/intro_screen.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';

import '../../widgets/button/button.dart';
import '../../widgets/input/input.dart';
import '../../widgets/text/text_app_bar.dart';

class ResetPasswordScreen extends StatelessWidget {
  const ResetPasswordScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Reset Password"),
        body: SingleChildScrollView(
          child: Container(
            padding: const EdgeInsets.all(16),
            child: const _ResetPasswordForm(),
          ),
        ));
  }
}

class _ResetPasswordForm extends StatefulWidget {
  const _ResetPasswordForm({Key? key}) : super(key: key);

  @override
  _ResetPasswordFormState createState() => _ResetPasswordFormState();
}

class _ResetPasswordFormState extends State<_ResetPasswordForm> {
  final TextEditingController _emailController = TextEditingController();

  OTPModel? _otpModel;

  bool isSubmitting = false;

  void _verifyOTP(String otpCode) async {
    if (_otpModel == null) {
      return;
    }
    if (!OTPModel.isValidCode(otpCode)) {
      Alert.message("Kode OTP tidak valid");
      return;
    }

    _otpModel!.code = otpCode;

    Alert.withLoading((done) async {
      var resetPasswordModel = await UserService.resetPasswordVerifyOTP(
          _emailController.text, _otpModel!);
      done();
      Screens.replace(_ResetPasswordScreen2(model: resetPasswordModel));
    });
  }

  void _handleResetPassword() async {
    String email = _emailController.text;
    if (email.isEmpty) {
      Alert.message("Email tidak boleh kosong");
      return;
    }
    Alert.withLoading((done) async {
      _otpModel = await UserService.resetPasswordOTP(email);
      done();
      Screens.to(BaseOTP(
          title: "Reset Password",
          subtitle: "Verifikasi",
          model: _otpModel!,
          description:
              "Masukkan 6 digit angka (OTP) yang kami telah kami kirimkan ke $email untuk melakukan reset password",
          onSubmit: _verifyOTP));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Input(
              controller: _emailController,
              label: "Alamat Email",
              hintText: "Masukkan Alamat Email",
              prefixIcon: const Icon(Icons.email),
            ),
            const SizedBox(height: 16),
            Button(
              width: double.infinity,
              onPressed: _handleResetPassword,
              isLoading: isSubmitting,
              child: const Text('Lanjutkan'),
            ),
          ],
        ),
      ],
    );
  }
}

class _ResetPasswordScreen2 extends StatefulWidget {
  const _ResetPasswordScreen2({Key? key, required this.model})
      : super(key: key);

  final ResetPasswordModel model;

  @override
  State<_ResetPasswordScreen2> createState() => _ResetPasswordState2();
}

class _ResetPasswordState2 extends State<_ResetPasswordScreen2> {
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmPasswordController =
      TextEditingController();

  void _onSubmit() async {
    if (passwordController.text.isEmpty) {
      Alert.message("Password baru tidak boleh kosong");
      return;
    }
    if (passwordController.text != confirmPasswordController.text) {
      Alert.message("Konfirmasi password tidak sesuai");
      return;
    }
    Alert.withLoading((done) async {
      var resp = await UserService.resetPassword(
          widget.model, passwordController.text);
      done();
      Alert.message(resp).then((_) {
        Screens.replaceAll(const IntroScreen());
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Buat Password Baru"),
      body: ListView(
        children: [
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Input(
                label: "Password Baru",
                controller: passwordController,
                isPassword: true,
              )),
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Input(
                label: "Konfirmasi Password Baru",
                controller: confirmPasswordController,
                isPassword: true,
              )),
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Button(
                width: double.infinity,
                onPressed: _onSubmit,
                child: const Text('Ganti Password'),
              )),
        ],
      ),
    );
  }
}
