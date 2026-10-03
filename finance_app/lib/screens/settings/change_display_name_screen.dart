import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class ChangeDisplayNameScreen extends StatelessWidget {
  const ChangeDisplayNameScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Ubah Nama Lengkap"),
      body: SafeArea(
        child: Container(
            padding: const EdgeInsets.all(16),
            child: const _ChangeDisplayNameContent()),
      ),
    );
  }
}

class _ChangeDisplayNameContent extends StatefulWidget {
  const _ChangeDisplayNameContent({Key? key}) : super(key: key);

  @override
  State<_ChangeDisplayNameContent> createState() =>
      _ChangeDisplayNameContentState();
}

class _ChangeDisplayNameContentState extends State<_ChangeDisplayNameContent> {
  FocusNode? focusNode = FocusNode();
  TextEditingController controller = TextEditingController();

  @override
  void initState() {
    super.initState();

    controller.text = AccountController.getInstance().authDetails!.displayName;

    Future.delayed(const Duration(milliseconds: 200), () {
      if (!mounted) return;
      if (focusNode!.hasFocus) return;
      if (focusNode == null) return;
      FocusScope.of(context).requestFocus(focusNode);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Input(
            focusNode: focusNode,
            label: 'Nama Lengkap',
            controller: controller),
        const SizedBox(height: 16),
        Button(
            onPressed: changeDisplayName,
            width: double.infinity,
            child: const Text('Ubah'))
      ],
    );
  }

  void changeDisplayName() {
    Alert.withLoading((done) async {
      await UserService.updateDisplayName(controller.text);
      done();
    });
  }

  @override
  void dispose() {
    // dispose the focus node
    focusNode!.dispose();
    focusNode = null;
    super.dispose();
  }
}
