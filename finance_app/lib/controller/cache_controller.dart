import 'package:finance_app/controller/controller.dart';

class CacheController extends Controller {
  final Map<String, dynamic> _caches = {};

  void set(String key, dynamic value) {
    _caches[key] = value;
    notifyListeners();
  }

  dynamic get(String key, dynamic value) {
    return _caches[key];
  }

  void clear() {
    _caches.clear();
    notifyListeners();
  }

  void remove(String key) {
    _caches.remove(key);
    notifyListeners();
  }
}
