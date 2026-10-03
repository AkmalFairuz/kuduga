import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class ChangePasswordScreen extends StatelessWidget {
  const ChangePasswordScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: TextAppBar("Perbarui Kata Sandi"),
      body: SingleChildScrollView(
        child: Padding(
          padding: EdgeInsets.all(16),
          child: _ChangePasswordContent(),
        ),
      ),
    );
  }
}

class _ChangePasswordContent extends StatefulWidget {
  const _ChangePasswordContent({Key? key}) : super(key: key);

  @override
  State<_ChangePasswordContent> createState() => _ChangePasswordContentState();
}

class _ChangePasswordContentState extends State<_ChangePasswordContent> {
  TextEditingController passwordController = TextEditingController();
  TextEditingController newPasswordController = TextEditingController();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Input(
          isPassword: true,
          label: 'Password saat ini',
          controller: passwordController,
        ),
        const SizedBox(height: 16),
        Input(
          isPassword: true,
          label: 'Password baru',
          controller: newPasswordController,
        ),
        const SizedBox(height: 8),
        const Text(
            'Password harus mengandung huruf dan angka\nPassword harus memiliki 8 karakter atau lebih\n\nSetelah mengganti password, perangkat lain yang menggunakan akun anda akan diminta untuk login ulang',
            style: TextStyle(color: Colors.grey)),
        const SizedBox(height: 16),
        Button(
            onPressed: changePassword,
            width: double.infinity,
            child: const Text('Ganti Password')),
      ],
    );
  }

  void changePassword() async {
    if (passwordController.text.isEmpty) {
      Alert.message("Password saat ini harus diisi");
      return;
    }
    if (newPasswordController.text.isEmpty) {
      Alert.message("Password baru harus diisi");
      return;
    }

    await Alert.withLoading((done) async {
      await UserService.updatePassword(
          passwordController.text, newPasswordController.text);
      Alert.message("Password berhasil diperbarui");
      done();
    });
  }
}
