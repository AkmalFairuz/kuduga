import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class OTPInput extends StatefulWidget {
  const OTPInput({
    Key? key,
    required this.onChanged,
    this.initialFocus = false,
  }) : super(key: key);

  final ValueChanged<String> onChanged;
  final bool initialFocus;

  @override
  State<StatefulWidget> createState() => _OTPInputState();
}

class _OTPInputState extends State<OTPInput> {
  final List<TextEditingController> _controllers = List.generate(
    6,
    (index) => TextEditingController(),
  );

  final List<FocusNode> _focusNodes = List.generate(
    6,
    (index) => FocusNode(),
  );

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildInput(context, 0),
        const SizedBox(width: 10),
        _buildInput(context, 1),
        const SizedBox(width: 10),
        _buildInput(context, 2),
        const SizedBox(width: 10),
        _buildInput(context, 3),
        const SizedBox(width: 10),
        _buildInput(context, 4),
        const SizedBox(width: 10),
        _buildInput(context, 5),
      ],
    );
  }

  @override
  void initState() {
    super.initState();
    if (widget.initialFocus) {
      Future.delayed(const Duration(milliseconds: 250)).then((_) {
        if (mounted && _controllers[0].text == '') {
          _focusNodes[0].requestFocus();
        }
      });
    }
  }

  Widget _buildInput(BuildContext context, int index) {
    var controller = _controllers[index];
    var focusNode = _focusNodes[index];
    return RawKeyboardListener(
        focusNode: FocusNode(),
        onKey: (event) {
          if (event.logicalKey == LogicalKeyboardKey.backspace && index > 0) {
            Future.delayed(const Duration(milliseconds: 10)).then((_) {
              if (mounted) {
                FocusScope.of(context).requestFocus(_focusNodes[index - 1]);
                _controllers[index - 1].clear();
              }
            });
            widget.onChanged(getValue());
            return;
          }
        },
        child: Container(
            width: 40,
            decoration: BoxDecoration(
              border: Border(
                bottom: BorderSide(
                  color: Colors.grey[400]!,
                ),
              ),
            ),
            child: TextField(
              focusNode: focusNode,
              textAlign: TextAlign.center,
              controller: controller,
              keyboardType: TextInputType.number,
              onChanged: (value) {
                handleChanged(value, focusNode, index, controller);
                widget.onChanged(getValue());
              },
            )));
  }

  void handleChanged(String value, FocusNode focusNode, int index,
      TextEditingController controller) {
    // detect if paste a six digit number, then fill all input
    if (value.length == 6) {
      for (int i = 0; i < 6; i++) {
        _controllers[i].text = value[i];
      }
      focusNode.unfocus();
      return;
    }

    // if user trying to input more than one character
    // then fill the input with the last character and move to next input
    if (value.length > 1) {
      controller.text = value[value.length - 1];
      if (index == 5) {
        focusNode.unfocus();
      } else {
        _focusNodes[index + 1].requestFocus();
      }
      return;
    }
    // if user input a character, then move to next input
    // if user input a character in the last input, then unfocus the input
    if (value.isNotEmpty) {
      if (index == 5) {
        focusNode.unfocus();
        return;
      }
      _focusNodes[index + 1].requestFocus();
      return;
    }
  }

  String getValue() {
    String value = "";
    for (int i = 0; i < 6; i++) {
      value += _controllers[i].text;
    }
    return value;
  }

  bool isAllFilled() {
    for (int i = 0; i < 6; i++) {
      if (_controllers[i].text.isEmpty) {
        return false;
      }
    }
    return true;
  }
}
