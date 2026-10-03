import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/screens/home/notifications_tab.dart';
import 'package:finance_app/screens/home/profile_tab.dart';
import 'package:finance_app/screens/home/purchase_history_tab.dart';
import 'package:finance_app/service/auth_service.dart';
import 'package:finance_app/service/notification_service.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/state/state_notifier.dart';
import 'package:finance_app/widgets/home/balance_card.dart';
import 'package:finance_app/widgets/home/carousel_info.dart';
import 'package:finance_app/widgets/home/services_card.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../utils/meta.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedTab = 0;
  int? _lastPop;

  List<Widget> _widgets = [];

  @override
  void initState() {
    _widgets = [
      const _HomeTab(),
      const PurchaseHistoryTab(),
      const NotificationsTab(),
      ProfileTab(onOpenPurchaseHistory: () {
        setState(() {
          _selectedTab = 1;
        });
      }),
    ];

    AuthService.extendToken();
    NotificationService.instance.pushToken();
    StateNotifier.getChannel("purchase")
        .listen(() => AccountController.getInstance().fetchBalance());
    UserService.updateDeviceUniqueID();

    super.initState();
  }

  Future<bool> onWillPop() async {
    if (_selectedTab != 0) {
      setState(() {
        _selectedTab = 0;
      });
      return false;
    }
    if (_lastPop == null || DateTime.now().millisecond - _lastPop! >= 1000) {
      _lastPop = DateTime.now().millisecond;
      return false;
    }
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        onWillPop: onWillPop,
        child: Scaffold(
          appBar: PreferredSize(
              preferredSize: const Size.fromHeight(0.0),
              child: AppBar(
                  backgroundColor:
                      Theme.of(context).appBarTheme.backgroundColor,
                  elevation: 0,
                  systemOverlayStyle: SystemUiOverlayStyle(
                    statusBarColor:
                        Theme.of(context).appBarTheme.backgroundColor,
                    statusBarIconBrightness: Brightness.light,
                    systemNavigationBarColor: Colors.grey[50],
                    systemNavigationBarIconBrightness: Brightness.dark,
                  ))),
          body: IndexedStack(
            index: _selectedTab,
            children: _widgets,
          ),
          bottomNavigationBar: BottomNavigationBar(
            iconSize: 24.sp,
            type: BottomNavigationBarType.fixed,
            selectedFontSize: 11.sp,
            unselectedFontSize: 11.sp,
            currentIndex: _selectedTab,
            onTap: (int index) {
              setState(() {
                _selectedTab = index;
              });
            },
            items: const [
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.home),
                label: "Beranda",
              ),
              BottomNavigationBarItem(
                icon: Icon(Icons.history),
                label: "Riwayat",
              ),
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.bell),
                label: "Notifikasi",
              ),
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.person),
                label: "Profil",
              ),
            ],
          ),
        ));
  }
}

class _HomeTab extends StatelessWidget {
  const _HomeTab();

  Future<void> _onRefresh() async {
    await AccountController.getInstance().fetchAll();
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
        onRefresh: _onRefresh,
        child: ListView(
          physics: const AlwaysScrollableScrollPhysics(),
          children: [
            Stack(
              children: [
                Container(
                  width: double.infinity,
                  height: 155,
                  decoration: BoxDecoration(
                    color: Meta.color[800]!,
                    borderRadius: const BorderRadius.only(
                      bottomLeft: Radius.elliptical(150, 40),
                      bottomRight: Radius.elliptical(150, 40),
                    ),
                  ),
                ),
                Container(
                  padding: const EdgeInsets.all(20),
                  child: const Column(
                    children: [BalanceCard()],
                  ),
                )
              ],
            ),
            const CarouselInfo(),
            const SizedBox(height: 24),
            const ServicesCard(),
          ],
        ));
  }
}
