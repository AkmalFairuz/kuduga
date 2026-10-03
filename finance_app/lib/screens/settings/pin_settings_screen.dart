import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/screens/auth/base_otp.dart';
import 'package:finance_app/screens/auth/pin_screen.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/button/list_button.dart';
import 'package:finance_app/widgets/button/list_button_view.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/text/modal_bottom_title.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

class PinSettingsScreen extends StatelessWidget {
  const PinSettingsScreen({super.key});

  Widget _buildBody(BuildContext context) {
    return AccountController.getInstance().authDetails!.hasPin
        ? ListButtonView(buttons: [
            ListButton(
                title: "Perbarui PIN",
                onPress: () {
                  Screens.to(_UpdatePinScreen());
                }),
            ListButton(
                title: "Hapus PIN",
                onPress: () {
                  ModalBottomSheet.show(context, (context) {
                    return Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const ModalBottomTitle("Hapus PIN"),
                        const SizedBox(height: 16),
                        const Text(
                            "Menghapus PIN dapat mengurangi keamanan akun anda. Apakah anda yakin ingin menghapus PIN?"),
                        const SizedBox(height: 16),
                        Row(
                          children: [
                            const Expanded(
                                child: Button(
                                    onPressed: Screens.back,
                                    variant: ButtonVariant.outline,
                                    child: Text("Batal"))),
                            const SizedBox(width: 16),
                            Expanded(
                                child: Button(
                                    onPressed: () {
                                      Screens.back(); // close modal

                                      Screens.to(PinScreen(
                                        title:
                                            "Masukkan PIN untuk menghapus PIN",
                                        onSubmit: (pin) async {
                                          await Alert.withLoading((done) async {
                                            var resp =
                                                await UserService.deletePin(
                                                    pin);
                                            done();
                                            AccountController.getInstance()
                                                .fetchAuthDetails();
                                            await Alert.message(resp);
                                            Screens.back();
                                          });
                                        },
                                      ));
                                    },
                                    child: const Text("Hapus PIN"))),
                          ],
                        )
                      ],
                    );
                  },
                      showDragHelper: false,
                      isDismissible: false,
                      enableDrag: false);
                }),
            ListButton(
                title: "Lupa PIN",
                onPress: () {
                  ModalBottomSheet.show(context, (context) {
                    return Column(
                      children: [
                        const ModalBottomTitle("Lupa PIN"),
                        const SizedBox(height: 16),
                        const Text(
                            "Untuk mengembalikan PIN, anda harus memasukkan kode OTP yang kami kirimkan ke email anda.\n\nSetelah itu, PIN anda akan terhapus dan anda dapat membuat PIN baru"),
                        const SizedBox(height: 16),
                        Button(
                            onPressed: () async {
                              Screens.back(); // close modal

                              Alert.withLoading((done) async {
                                final otpModel =
                                    await UserService.recoveryPinOTP();
                                done();

                                Screens.to(BaseOTP(
                                    model: otpModel,
                                    onSubmit: (otpCode) async {
                                      otpModel.code = otpCode;

                                      Alert.withLoading((done) async {
                                        String resp =
                                            await UserService.recoveryPin(
                                                otpModel);
                                        done();
                                        AccountController.getInstance()
                                            .fetchAuthDetails();
                                        await Alert.message(resp);
                                        Screens.back();
                                      });
                                    },
                                    title: "Lupa PIN",
                                    subtitle: "Lupa PIN",
                                    description:
                                        "Masukkan kode OTP yang kami kirimkan ke email anda untuk menghapus PIN"));
                              });
                            },
                            width: double.infinity,
                            child: const Text("Lanjutkan"))
                      ],
                    );
                  }, showDragHelper: false);
                }),
          ])
        : ListView(
            children: [
              const SizedBox(height: 32),
              const Text(
                "Anda belum menyetel PIN",
                textAlign: TextAlign.center,
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Expanded(
                      child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          child: const Text(
                            "Disarankan untuk menyetel PIN untuk mengurangi resiko akun terkena hack.\n\nSetelah menyetel PIN, anda akan diminta untuk memasukkan PIN setiap kali membuka aplikasi dan transaksi.",
                            textAlign: TextAlign.center,
                          )))
                ],
              ),
              const SizedBox(height: 32),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Button(
                      onPressed: () {
                        Screens.to(_CreatePinScreen());
                      },
                      child: const Text("Buat PIN"))
                ],
              )
            ],
          );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: const TextAppBar("Pengaturan PIN"),
        body: Dyn(() => _buildBody(context), AccountController.getInstance()));
  }
}

class _CreatePinScreen extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _CreatePinState();
}

class _CreatePinState extends State<_CreatePinScreen> {
  final TextEditingController pinController = TextEditingController();
  final TextEditingController confirmPinController = TextEditingController();

  void _onSubmit() async {
    String pin = pinController.text;
    if (!Meta.isStringOnlyNumbers(pin)) {
      Alert.message("PIN harus berupa angka");
      return;
    }
    if (pin.length != 6) {
      Alert.message("PIN harus 6 digit");
      return;
    }
    if (confirmPinController.text != pin) {
      Alert.message("Konfirmasi PIN tidak sesuai");
      return;
    }

    Alert.withLoading((done) async {
      var resp = await UserService.createPin(pin);
      done();
      AccountController.getInstance().fetchAuthDetails();
      await Alert.message(resp);
      if (mounted) {
        Navigator.of(context).pop(true);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Buat PIN"),
      body: ListView(
        children: [
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Input(
              controller: pinController,
              label: "PIN baru (6 digit)",
              hintText: "Masukkan PIN baru",
              keyboardType: TextInputType.number,
              isPassword: true,
            ),
          ),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Input(
              controller: confirmPinController,
              label: "Konfirmasi PIN baru (6 digit)",
              hintText: "Masukkan PIN baru ulang",
              keyboardType: TextInputType.number,
              isPassword: true,
            ),
          ),
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child:
                  Button(onPressed: _onSubmit, child: const Text("Buat PIN")))
        ],
      ),
    );
  }
}

class _UpdatePinScreen extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _UpdatePinState();
}

class _UpdatePinState extends State<_UpdatePinScreen> {
  final TextEditingController pinController = TextEditingController();
  final TextEditingController newPinController = TextEditingController();
  final TextEditingController confirmNewPinController = TextEditingController();

  void _onSubmit() async {
    String pin = pinController.text;
    if (!Meta.isPin(pin)) {
      Alert.show(title: "Kesalahan", message: "PIN salah");
      return;
    }

    String newPin = newPinController.text;
    if (!Meta.isStringOnlyNumbers(newPin)) {
      Alert.message("PIN baru harus berupa angka");
      return;
    }
    if (newPin.length != 6) {
      Alert.message("PIN baru harus 6 digit");
      return;
    }
    if (confirmNewPinController.text != newPin) {
      Alert.message("Konfirmasi PIN baru tidak sesuai");
      return;
    }

    try {
      Alert.withLoading((done) async {
        var resp = await UserService.updatePin(pin, newPin);
        done();
        await Alert.show(title: "Berhasil", message: resp);
        if (mounted) {
          Navigator.of(context).pop(true);
        }
      });
    } finally {
      AccountController.getInstance().fetchAuthDetails();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Perbarui PIN"),
      body: ListView(
        children: [
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Input(
              controller: pinController,
              label: "PIN saat ini (6 digit)",
              hintText: "Masukkan PIN saat ini",
              keyboardType: TextInputType.number,
              isPassword: true,
            ),
          ),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Input(
              controller: newPinController,
              label: "PIN baru (6 digit)",
              hintText: "Masukkan PIN baru",
              keyboardType: TextInputType.number,
              isPassword: true,
            ),
          ),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Input(
              controller: confirmNewPinController,
              label: "Konfirmasi PIN baru (6 digit)",
              hintText: "Masukkan PIN baru ulang",
              keyboardType: TextInputType.number,
              isPassword: true,
            ),
          ),
          const SizedBox(height: 16),
          Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Button(
                  onPressed: _onSubmit, child: const Text("Perbarui PIN")))
        ],
      ),
    );
  }
}
