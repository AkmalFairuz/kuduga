import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'input.dart';

// CurrencyInput is a widget that is used to input currency.
// It uses Input widget as its base.
// when typing the input, it will automatically format the input to currency format.
// it add points every 3 digits from the right.

class CurrencyInput extends StatefulWidget {
  const CurrencyInput(
      {Key? key,
      this.hintText,
      this.controller,
      this.error,
      this.label,
      this.onChanged})
      : super(key: key);

  final String? hintText;
  final TextEditingController? controller;
  final String? error;
  final String? label;
  final Function(int)? onChanged;

  @override
  State<CurrencyInput> createState() => _CurrencyInputState();
}

class _CurrencyInputState extends State<CurrencyInput> {
  final TextEditingController _controller = TextEditingController();

  @override
  void initState() {
    super.initState();
    _controller.text = widget.controller?.text ?? '';
  }

  @override
  Widget build(BuildContext context) {
    return Input(
      label: widget.label,
      error: widget.error,
      hintText: widget.hintText,
      keyboardType: TextInputType.number,
      controller: _controller,
      prefixIcon: Row(
        mainAxisSize: MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(' Rp',
              style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                  color: Colors.grey[700]!)),
        ],
      ),
      onChanged: (value) {
        final text = value.replaceAll('.', '');
        final number = int.tryParse(text);
        if (number != null) {
          final currency = number.toString().replaceAllMapped(
                RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
                (match) => '${match[1]}.',
              );
          _controller.value = TextEditingValue(
            text: currency,
            selection: TextSelection.collapsed(offset: currency.length),
          );
          widget.controller?.text = number.toString();
          widget.onChanged?.call(number);
        }
      },
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
        LengthLimitingTextInputFormatter(7)
      ],
    );
  }
}
