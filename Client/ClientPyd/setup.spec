# -*- mode: python ; coding: utf-8 -*-

block_cipher = None

a = Analysis(['Virus.py'],
             pathex=[],
             binaries=[],
             datas=[],
             hiddenimports=[
                            'os',
                            'os.path',
                            'pathlib',
                            'io',
                            'time',
                            'json',
                            'base64',
                            'strings',
                            'github.com/gin-gonic/gin',
                            'net',
                            'log',
                            'http',
                            'exec',
                            'fmt',
                            'filepath',
                            # 注意最后需要逗号
                            ],
             hookspath=[],  # 此处注意语法规范
             runtime_hooks=[],
             excludes=[],
             win_no_prefer_redirects=False,
             win_private_assemblies=False,
             cipher=block_cipher,
             noarchive=False)

pyz = PYZ(a.pure, a.zipped_data,
          cipher=block_cipher)

exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          a.datas,
          [],
          name='Virus',
          debug=False,
          bootloader_ignore_signals=False,
          strip=False,
          upx=True,
          upx_exclude=[],
          runtime_tmpdir=None,
          console=True)
