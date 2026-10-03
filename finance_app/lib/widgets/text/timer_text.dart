import 'dart:async';

import 'package:flutter/cupertino.dart';

class TimerText extends StatefulWidget {
  const TimerText({Key? key, required this.duration, this.onFinish, this.style})
      : super(key: key);

  final Duration duration;
  final VoidCallback? onFinish;
  final TextStyle? style;

  @override
  State<TimerText> createState() => _TimerTextState();
}

class _TimerTextState extends State<TimerText> {
  late Duration _duration;
  late Timer _timer;

  @override
  void initState() {
    _duration = widget.duration;
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_duration.inSeconds == 0) {
        _timer.cancel();
        if (widget.onFinish != null) {
          widget.onFinish!();
        }
      } else {
        setState(() {
          _duration = _duration - const Duration(seconds: 1);
        });
      }
    });
    super.initState();
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  String _format(Duration duration) {
    // format HH:MM:SS
    String twoDigits(int n) => n.toString().padLeft(2, "0");
    String twoDigitMinutes = twoDigits(duration.inMinutes.remainder(60));
    String twoDigitSeconds = twoDigits(duration.inSeconds.remainder(60));
    return "${twoDigits(duration.inHours)}:$twoDigitMinutes:$twoDigitSeconds";
  }

  @override
  Widget build(BuildContext context) {
    return Text(_format(_duration), style: widget.style);
  }
}
