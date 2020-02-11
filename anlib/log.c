#include <android/log.h>

void anlib_logs(int fs,const char*tag,const char*conts){
    if(fs==1){
        __android_log_print(ANDROID_LOG_INFO,tag,conts);
    }else if(fs==2){
        __android_log_print(ANDROID_LOG_DEBUG,tag,conts);
    }else if(fs==3){
        __android_log_print(ANDROID_LOG_ERROR,tag,conts);
    }
}