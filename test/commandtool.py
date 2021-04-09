# !/usr/bin/python
# -*- coding: utf-8 -*-
import argparse
import ctypes
import os

def getfile():
    """
    获取动态库文件
    :return: 返回获取的文件
    """
    filepath = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    file = ctypes.CDLL(filepath + "/matchCommand.so")
    return file

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    # 添加参数
    parser.add_argument("-m", nargs="+", help="Explanation of all parameters of single commands", dest="m", type=str)
    parser.add_argument("-ml", nargs="+", help="Explanation of all parameters of a multiple command", dest="ml", type=str)
    parser.add_argument("-e", nargs="+", help="Interpretation of current parameters of a single command", dest="e", type=str)
    parser.add_argument("-el", nargs="+", help="Interpretation of current parameters of multiple commands", dest="el", type=str)
    parser.add_argument("-o", nargs="?", help="operating system", dest="o", type=str)
    args = parser.parse_args()

    if args.m:
        # 导入.so文件，获取man_cmd方法（路径可根据需求进行更改）
        manCmd = getfile().man
        # 设置期望的传参类型
        manCmd.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # 给定操作系统类型的值
        if args.o:
            os_type = "".join(args.o)
        else:
            os_type = ""
        manCmd.restype = ctypes.c_char_p
        # 以utf-8编码对os_type进行编码，获得bytes类型对象
        os_type = os_type.encode("utf-8")
        # 给定要查询含义的命令行
        cmd = " ".join(args.m)
        # 以utf-8编码对cmd进行编码，获得bytes类型对象
        cmd = cmd.encode("utf-8")
        # 调用man方法
        result = manCmd(cmd, os_type)
        print(result.decode("utf-8"))
    elif args.ml:
        # 导入.so文件，获取man_cmd_list方法（路径可根据需求进行更改）
        manCmdList = getfile().man_list
        # 设置期望的传参类型
        manCmdList.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # 给定操作系统类型的值
        if args.o:
            os_type = "".join(args.o)
        else:
            os_type=""
        manCmdList.restype = ctypes.c_char_p
        # 以utf-8编码对os_type进行编码，获得bytes类型对象
        os_type = os_type.encode("utf-8")
        # 给定要查询含义的命令行
        cmd = " ".join(args.ml)
        # 以utf-8编码对cmd进行编码，获得bytes类型对象
        cmd = cmd.encode("utf-8")
        # 调用man_list方法
        result = manCmdList(cmd, os_type)
        print(result.decode("utf-8"))
    elif args.e:
        explainCmd = getfile().explain_cmd
        # 设置期望的传参类型
        explainCmd.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # 给定操作系统类型的值
        if args.o:
            os_type = "".join(args.o)
        else:
            os_type = ""
        explainCmd.restype = ctypes.c_char_p
        # 以utf-8编码对os_type进行编码，获得bytes类型对象
        os_type = os_type.encode("utf-8")
        # 给定要查询含义的命令行
        cmd = " ".join(args.e)
        # 以utf-8编码对cmd进行编码，获得bytes类型对象
        cmd = cmd.encode("utf-8")
        # 调用explain_cmd方法
        result = explainCmd(cmd, os_type)
        print(result.decode("utf-8"))
    elif args.el:
        explainCmdList = getfile().explain_cmd_list
        # 设置期望的传参类型
        explainCmdList.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # 给定操作系统类型的值s
        if args.o:
            os_type = "".join(args.o)
        else:
            os_type = ""
        explainCmdList.restype = ctypes.c_char_p
        # 以gb2312编码对os_type进行编码，获得bytes类型对象
        os_type = os_type.encode("utf-8")
        # 给定要查询含义的命令行
        cmd = " ".join(args.el)
        # 以utf-8编码对cmd进行编码，获得bytes类型对象
        cmd = cmd.encode("utf-8")
        # 调用explain_cmd_list方法
        result = explainCmdList(cmd, os_type)
        print(result.decode("utf-8"))
