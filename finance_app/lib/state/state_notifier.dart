import 'package:flutter/cupertino.dart';

class StateNotifier {
  static final Map<String, StateNotifierChannel> _channels = {};

  static StateNotifierChannel createChannel(String channelName) {
    if (!_channels.containsKey(channelName)) {
      _channels[channelName] = StateNotifierChannel();
    }
    return _channels[channelName]!;
  }

  static StateNotifierChannel getChannel(String channelName) {
    var ret = _channels[channelName];
    ret ??= createChannel(channelName);
    return ret;
  }

  static Widget listenWidget(
      String channelName, Widget Function(BuildContext context) onNotified) {
    createChannel(channelName);
    return ListenableBuilder(
        listenable: _channels[channelName]!,
        builder: (context, _) {
          return onNotified(context);
        });
  }

  static void notify(String channelName) {
    if (_channels.containsKey(channelName)) {
      _channels[channelName]!.notify();
    }
  }
}

class StateNotifierChannel with ChangeNotifier {
  void notify() {
    notifyListeners();
  }

  void listen(Function() onNotified) {
    addListener(onNotified);
  }

  void closeListener(Function() onNotified) {
    removeListener(onNotified);
  }
}
