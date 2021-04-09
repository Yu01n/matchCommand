# !/usr/bin/python
# -*- coding: utf-8 -*-
import ctypes
import os
import time

import xlrd

import pandas as pd


def getfile():
    """
    获取动态库文件
    :return: 返回获取的文件
    """
    filepath = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    file = ctypes.CDLL(filepath + "/matchCommand.so")
    return file

def man(cmd, os_type=""):
    """
    输出单个命令行的全部参数信息，如需要查询多个，以";"分隔
    :param cmd: 单个命令行
    :param os_type: 操作系统类型，可传linux和windows，也可以不传值
    :return: 返回json格式字符串
    """
    # 设置环境变量值
    filepath = os.getenv('getfilepath', '..')
    # 导入.so文件，获取man_list方法（路径可根据需求进行更改）
    manCmd = getfile().man
    # 设置期望的传参类型和返回类型
    manCmd.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
    manCmd.restype = ctypes.c_char_p
    # 给定获取json文件路径的值,以utf-8编码对filepath进行编码，获得bytes类型对象
    filepath = filepath.encode("utf-8")
    # 给定操作系统类型的值,以utf-8编码对os_type进行编码，获得bytes类型对象
    os_type = os_type.encode("utf-8")
    # 给定要查询含义的命令行,以utf-8编码对cmd进行编码，获得bytes类型对象
    cmd = cmd.encode("utf-8")
    # 调用man方法
    result = manCmd(filepath, cmd, os_type)
    return result.decode("utf-8")


def man_list(cmd, os_type=""):
    """
    输出多个命令行的全部参数信息，如需要查询多个，以";"分隔
    :param cmd: 多个命令行
    :param os_type: 操作系统类型，可传linux和windows，也可以不传值
    :return: 返回json格式字符串
    """
    # 设置环境变量值
    filepath = os.getenv('getfilepath', '..')
    # 导入.so文件，获取man方法（路径可根据需求进行更改）
    manCmdList = getfile().man_list
    # 设置期望的传参类型和返回类型
    manCmdList.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
    manCmdList.restype = ctypes.c_char_p
    # 给定获取json文件路径的值,以utf-8编码对filepath进行编码，获得bytes类型对象
    filepath = filepath.encode("utf-8")
    # 给定操作系统类型的值, 以utf-8编码对os_type进行编码，获得bytes类型对象
    os_type = os_type.encode("utf-8")
    # 给定要查询含义的命令行,以utf-8编码对cmd进行编码，获得bytes类型对象
    cmd = cmd.encode("utf-8")
    # 调用man_list方法
    result = manCmdList(filepath, cmd, os_type)
    return result.decode("utf-8")


def explain_cmd(cmd, os_type=""):
    """
    输出单个命令行的当前参数信息
    :param cmd: 单个命令行
    :param os_type: 操作系统类型，可传linux和windows，也可以不传值
    :return: 返回json格式字符串
    """
    # 设置环境变量值
    filepath = os.getenv('getfilepath', '..')
    explainCmd = getfile().explain_cmd
    # 设置期望的传参类型和返回类型
    explainCmd.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
    explainCmd.restype = ctypes.c_char_p
    # 给定获取json文件路径的值,以utf-8编码对filepath进行编码，获得bytes类型对象
    filepath = filepath.encode("utf-8")
    # 给定操作系统类型的值,以utf-8编码对os_type进行编码，获得bytes类型对象
    os_type = os_type.encode("utf-8")
    # 给定要查询含义的命令行,以utf-8编码对cmd进行编码，获得bytes类型对象
    cmd = cmd.encode("utf-8")
    # 调用explain_cmd方法
    result = explainCmd(filepath, cmd, os_type)
    return result.decode("utf-8")


def explain_cmd_list(cmd, os_type=""):
    """
    输出多个命令行的当前参数信息，如需要查询多个，以";"分隔
    :param cmd: 多个命令行
    :param os_type: 操作系统类型，可传linux和windows，也可以不传值
    :return: 返回json格式字符串
    """
    # 设置环境变量值
    filepath = os.getenv('getfilepath', '..')
    explainCmdList = getfile().explain_cmd_list
    # 设置期望的传参类型和返回类型
    explainCmdList.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
    explainCmdList.restype = ctypes.c_char_p
    # 给定获取json文件路径的值,以utf-8编码对filepath进行编码，获得bytes类型对象
    filepath = filepath.encode("utf-8")
    # 给定操作系统类型的值, 以utf-8编码对os_type进行编码，获得bytes类型对象
    os_type = os_type.encode("utf-8")
    # 给定要查询含义的命令行,以utf-8编码对cmd进行编码，获得bytes类型对象
    cmd = cmd.encode("utf-8")
    # 调用explain_cmd_list方法
    result = explainCmdList(filepath, cmd, os_type)
    return result.decode("utf-8")


if __name__ == '__main__':

    # time_start=time.time()
    # for i in range(1,10000):
    #     cmdinfo1 = man("rpm -qf /usr/bin/mongod")
    # print(cmdinfo1)
    # time_end=time.time()
    # print('time cost', time_end-time_start, 's')

    # x1 = xlrd.open_workbook("D:\matchCommand\cmd.xlsx")
    # table = x1.sheet_by_index(0)
    # nrows = table.nrows
    # print("表格一共有", nrows, "行")
    #
    # time_start=time.time()
    # for i in range (0, nrows):
    #     s = table.row_values(i)
    #     cmd_info = explain_cmd_list(str(s))
    #     print(cmd_info)
    # time_end = time.time()
    # print('time cost', time_end-time_start, 's')


    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("")
    print(cmdinfo4)
    cmdinfo4 = explain_cmd("rm -rf test.txt")
    print(cmdinfo4)









