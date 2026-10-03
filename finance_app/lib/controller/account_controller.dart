import 'package:finance_app/controller/controller.dart';
import 'package:finance_app/model/auth_details.dart';
import 'package:finance_app/service/auth_service.dart';
import 'package:finance_app/service/balance_service.dart';

class AccountController extends Controller {
  static AccountController? _instance;

  static void destroyInstance() {
    _instance = null;
  }

  static AccountController getInstance() {
    if (_instance == null) {
      AccountController c = AccountController();
      c.fetchAll();
      _instance = c;
    }
    return _instance!;
  }

  int? _balance = 0;
  AuthDetailsModel? _authDetails;
  bool fetched = false;

  int? get balance => _balance;
  AuthDetailsModel? get authDetails => _authDetails;
  bool get isFetched => fetched;

  Future<void> fetchAll() async {
    _authDetails = await AuthService.getAuthDetails();
    _balance = await BalanceService.getBalance();
    fetched = true;
    notifyListeners();
  }

  Future<void> fetchAuthDetails() async {
    _authDetails = await AuthService.getAuthDetails();
    notifyListeners();
  }

  Future<void> fetchBalance() async {
    _balance = await BalanceService.getBalance();
    notifyListeners();
  }
}
