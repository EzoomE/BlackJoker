from setuptools import setup
from Cython.Build import cythonize

setup(
    name='StartWin32Joker',
    ext_modules=cythonize("main.pyx"),
)
