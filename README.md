# SnowRunnerWSGConverter

The purpose of this project is to help transfer/convert the XGP version of SnowRunner's save files to their normal save file format that is found with Steam and EGS versions.

The current binary is in the `dist` dir, but you can freely build a new one yourself if you have Go installed by running `./dist/build.bat`

# Sample Usage

Available options:

	-dest string
		The destination directory for converted save files (default ".")
	-ext string
		File extension to save all converted files with (default "cfg")
	-list-only
		List WSG save file mappings only
	-src string
		The source save directory containing WSG save files (default ".")

---
Converting files:

	$ SnowRunnerWSGConverter.exe --src=C:/Path/To/Src --dest=C:/Path/To/Dest   
	Using 'C:/Path/To/Src' as source directory
	Using 'C:/Path/To/Dest' as destination directory

	Found container file: container.13
	Converting WSG save files..

	Copied: C:/Path/To/Src/2E58A7E153FA42EEA6EB1DFCB7ACBCE8 -> C:/Path/To/Dest/CommonSslSave.cfg
	Copied: C:/Path/To/Src/C2C8F614EECE48C9A42C70C0EBBC653D -> C:/Path/To/Dest/fog_level_us_01_01.cfg
	Copied: C:/Path/To/Src/B688C916B72A4824A6D743BCEAE4BAAC -> C:/Path/To/Dest/achievements.cfg
	Copied: C:/Path/To/Src/35156D76DBF34DFBB4BCE3A2CC7E4638 -> C:/Path/To/Dest/CompleteSave.cfg
	Copied: C:/Path/To/Src/685741B717044FE1973F94E47B40DA21 -> C:/Path/To/Dest/user_settings.cfg
	Copied: C:/Path/To/Src/931363F3147A4D81B86556DD9BD1D140 -> C:/Path/To/Dest/user_social_data.cfg
	Copied: C:/Path/To/Src/4E5442B00C2B4AF8A1F44C04AD92E16C -> C:/Path/To/Dest/user_profile.cfg
	Copied: C:/Path/To/Src/33106F9817A74200817EE7A8DACA8D22 -> C:/Path/To/Dest/sts_level_us_01_01.cfg
	
---

Listing file mappings:

	$ SnowRunnerWSGConverter.exe --list-only --src=C:/Path/To/Src --dest=C:/Path/To/Dest       
	Using 'C:/Path/To/Src' as source directory

	Found container file: container.13
	WSG save file mapping:

	CommonSslSave -> 2E58A7E153FA42EEA6EB1DFCB7ACBCE8
	CompleteSave -> 35156D76DBF34DFBB4BCE3A2CC7E4638
	achievements -> B688C916B72A4824A6D743BCEAE4BAAC
	fog_level_ru_02_02 -> DEEF67BF193B4161AC945CFCCB898C46
	fog_level_us_01_01 -> C2C8F614EECE48C9A42C70C0EBBC653D
	fog_level_us_01_02 -> 6F52B133BE274115A129C6FC84BEA368
	fog_level_us_02_01 -> 13EF26E45354400D9D254D0937A97C0B
	sts_level_ru_02_02 -> B6EF8F0949944E45B21F74CB260CBBD9
	sts_level_us_01_01 -> 33106F9817A74200817EE7A8DACA8D22
	sts_level_us_01_02 -> 46857E7972564AEE884817C73BF77995
	sts_level_us_02_01 -> 9BAE160CCBD749DB9BF236CBAF9F28B0
	user_profile -> 4E5442B00C2B4AF8A1F44C04AD92E16C
	user_settings -> 685741B717044FE1973F94E47B40DA21
	user_social_data -> 931363F3147A4D81B86556DD9BD1D140
