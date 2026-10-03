import 'package:finance_app/model/app_layout_info.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/image/cached_svg.dart';
import 'package:finance_app/widgets/image/image.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../../service/app_service.dart';

class ServicesCard extends StatefulWidget {
  const ServicesCard({super.key});

  @override
  State<ServicesCard> createState() => _ServicesCardState();
}

class _ServicesCardState extends State<ServicesCard> {
  AppLayoutInfo? _appLayoutInfo;

  @override
  void initState() {
    super.initState();

    AppService.getLayoutInfo().then((value) {
      setState(() {
        _appLayoutInfo = value;
      });
    });
  }

  Widget _buildServices() {
    List<Widget> children = [];
    for (final svc in _appLayoutInfo!.services) {
      children.add(Padding(
        padding: EdgeInsets.symmetric(horizontal: 20.sp),
        child: Text(
          svc.title,
          style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 17.sp,
              color: Theme.of(context).textTheme.bodyLarge!.color),
        ),
      ));
      children.add(SizedBox(height: 16.sp));
      children.add(GridView.count(
        physics: const NeverScrollableScrollPhysics(),
        crossAxisCount: 4,
        mainAxisSpacing: 16.sp,
        shrinkWrap: true,
        children: svc.services.map((e) => ServiceBox(serviceInfo: e)).toList(),
      ));
      children.add(SizedBox(height: 10.sp));
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: children,
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_appLayoutInfo == null) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Padding(
              padding: EdgeInsets.symmetric(horizontal: 20),
              child: Skeleton(width: 160, height: 20)),
          const SizedBox(height: 32),
          GridView.count(
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 4,
            mainAxisSpacing: 16,
            shrinkWrap: true,
            children: const [
              ServiceBox(),
              ServiceBox(),
              ServiceBox(),
              ServiceBox(),
              ServiceBox(),
              ServiceBox(),
            ],
          )
        ],
      );
    }
    return _buildServices();
  }
}

class ServiceBox extends StatelessWidget {
  const ServiceBox({Key? key, this.serviceInfo}) : super(key: key);

  final ServiceInfo? serviceInfo;

  Widget _buildIcon(BuildContext context) {
    Color color = ColorHelper.theme(context,
        dark: Meta.color[400], light: Meta.color[600])!;
    switch (serviceInfo!.iconType) {
      case "svg":
        return Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CachedSvg(
              serviceInfo!.iconData,
              width: 45.sp,
            )
          ],
        );
      case "image":
        return ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: XImage(serviceInfo!.iconData, width: 38.sp, height: 38.sp));
      case "material":
        return Icon(
            IconData(int.parse(serviceInfo!.iconData),
                fontFamily: 'MaterialIcons'),
            color: color,
            size: 38.sp);
      case "cupertino":
        return Icon(
            IconData(int.parse(serviceInfo!.iconData),
                fontFamily: CupertinoIcons.iconFont,
                fontPackage: CupertinoIcons.iconFontPackage),
            color: color,
            size: 38.sp);
    }
    return Icon(Icons.question_mark, size: 38.sp, color: color);
  }

  Future<void> _handleTap() async {
    if (serviceInfo == null) {
      return;
    }
    await serviceInfo!.doAction();
  }

  @override
  Widget build(BuildContext context) {
    bool isSkeleton = serviceInfo == null;
    return GestureDetector(
        onTap: serviceInfo != null ? _handleTap : null,
        child: Column(
          children: [
            isSkeleton
                ? Skeleton(
                    width: 54.sp,
                    height: 54.sp,
                    decoration:
                        BoxDecoration(borderRadius: BorderRadius.circular(20)),
                  )
                : Container(
                    padding: const EdgeInsets.all(2),
                    width: 54.sp,
                    height: 54.sp,
                    decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(20.sp),
                        gradient: LinearGradient(
                            begin: Alignment.topLeft,
                            end: Alignment.bottomRight,
                            colors: [
                              ColorHelper.theme(context,
                                  dark: Meta.color[100],
                                  light: Meta.color[50])!,
                              ColorHelper.theme(context,
                                  dark: Meta.color[300],
                                  light: Meta.color[100])!,
                            ])),
                    child: _buildIcon(context),
                  ),
            SizedBox(
              height: 12.sp,
            ),
            Container(
                padding: const EdgeInsets.symmetric(horizontal: 4),
                child: isSkeleton
                    ? Skeleton(height: 10.sp, width: 40.sp)
                    : Text(
                        serviceInfo!.name,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            color: ColorHelper.theme(context,
                                dark: Meta.color[300]!,
                                light: Meta.color[700]!),
                            fontWeight: FontWeight.bold,
                            fontSize: 11.sp),
                      ))
          ],
        ));
  }
}
