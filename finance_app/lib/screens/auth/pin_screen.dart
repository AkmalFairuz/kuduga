import 'package:finance_app/utils/meta.dart';
import 'package:flutter/material.dart';

class PinScreen extends StatefulWidget {
  const PinScreen(
      {super.key, this.onSubmit, required this.title, this.hasLogout = false});

  final Function(String)? onSubmit;
  final String title;
  final bool hasLogout;

  @override
  State<StatefulWidget> createState() => _PinScreenState();
}

class _PinScreenState extends State<PinScreen> {
  String pin = "";

  void handleChange(int num) {
    if (pin.length >= 6) {
      return;
    }
    if (pin.length >= 5) {
      pin += num.toString();
      if (widget.onSubmit != null) {
        widget.onSubmit!(pin);
      } else {
        Navigator.of(context).pop(pin);
      }
      setState(() {
        pin = "";
      });
    } else {
      setState(() {
        pin += num.toString();
      });
    }
  }

  List<Widget> buildPinWidgets() {
    List<Widget> ret = [];
    for (int i = 1; i <= 9; i++) {
      ret.add(_PinButton(
          num: i,
          onPress: () {
            handleChange(i);
          }));
    }
    ret.add(Container());
    ret.add(_PinButton(
        num: 0,
        onPress: () {
          handleChange(0);
        }));
    ret.add(_PinButton(
        num: -1,
        onPress: () {
          if (pin.isEmpty) {
            return;
          }
          setState(() {
            pin = pin.substring(0, pin.length - 1);
          });
        }));
    return ret;
  }

  Widget buildPinIndicator() {
    List<Widget> children = [];
    for (int i = 1; i <= 6; i++) {
      children.add(Container(
        margin: const EdgeInsets.symmetric(horizontal: 10),
        decoration: BoxDecoration(
            color: pin.length >= i ? Colors.grey[100] : Meta.color[500],
            borderRadius: BorderRadius.circular(999)),
        width: 16,
        height: 16,
      ));
    }
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: children,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Meta.color[800],
        body: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Text(
                widget.title,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: Colors.grey[50],
                    fontWeight: FontWeight.bold,
                    fontSize: 18),
              ),
            ),
            const SizedBox(height: 50),
            buildPinIndicator(),
            const SizedBox(height: 25),
            if (widget.hasLogout)
              TextButton(
                  onPressed: () {},
                  child: Text(
                    "Lupa PIN?",
                    style: TextStyle(
                        color: Colors.grey[50], fontWeight: FontWeight.bold),
                  )),
            const SizedBox(height: 5),
            Container(
              constraints: const BoxConstraints(
                maxWidth: 340,
              ),
              child: GridView.count(
                padding: const EdgeInsets.all(16),
                mainAxisSpacing: 8,
                crossAxisSpacing: 8,
                crossAxisCount: 3,
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                children: buildPinWidgets(),
              ),
            ),
          ],
        ));
  }
}

class _PinButton extends StatelessWidget {
  const _PinButton({Key? key, required this.num, required this.onPress})
      : super(key: key);

  final int num;
  final VoidCallback onPress;

  @override
  Widget build(BuildContext context) {
    return InkWell(
        onTap: onPress,
        borderRadius: BorderRadius.circular(999),
        child: Container(
            padding: const EdgeInsets.all(8),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                num == -1
                    ? Icon(Icons.backspace_rounded, color: Colors.grey[50])
                    : Text(
                        num.toString(),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: Colors.grey[50],
                            fontSize: 28),
                      )
              ],
            )));
  }
}
