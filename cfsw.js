const PASSWORD = "admin";
const DOWNLAODURL = "https://cdn.truefilesize.com/test/test-50mb.bin";
export default {
    async fetch(request, env, ctx) {
        const MY_KV = env.KV;
        const method = request.method;
        let 请求URL文本 = request.url.replace(/%5[Cc]/g, '').replace(/\\/g, '');
        const url = new URL(请求URL文本);
        const cookies = request.headers.get('cookie') || 'null';
        const Cookie = getCookie(cookies,'auth');
        const salt ='a7F3k9Qx2Lm8Zp4Na7F3k9Qx2Lm8Zp4N';
        const gKEY = await MD5MD5(PASSWORD+salt);
        const urlrouter = url.pathname.slice(1).toLowerCase();
        if (urlrouter === PASSWORD){
            const resp = new Response(JSON.stringify({ success: gKEY }), { status: 200, headers: { 'Content-Type': 'application/json;charset=utf-8' } });
            resp.headers.set('Set-Cookie', `auth=${gKEY}; Path=/; Max-Age=86400; HttpOnly; Secure; SameSite=Lax`);
            return resp;
        }else if (urlrouter === 'download'){
            const response = await fetch(DOWNLAODURL,{
                headers:{
                    "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36",
                    "cache-control": "no-cache",
                }
            });
            const headers = new Headers();

            for (const key of [
                "content-type",
                "content-length",
                "content-range",
                "accept-ranges"
            ]) {
                const value = response.headers.get(key);
                if (value) {
                    headers.set(key, value);
                }
            }
            headers.set(
                "Content-Disposition",
                'attachment; filename="test-100mb.bin"'
            );
            if (!response.ok) {
                return new Response("download failed", {
                    status: 500
                });
            }
            return new Response(response.body, {
                status: 200,
                headers
            });
        }else if (Cookie === gKEY){
            if (urlrouter === 'ip'){
                if (method === 'GET'){
                    let list = await env.KV.get("ip", {
                        type: "json"
                    }) || [];
                    return new Response(JSON.stringify(list),{
                        status: 200,
                        headers: {
                            'Content-Type': 'application/json;charset=utf-8'
                        }
                    })
                }else if (method === 'POST'){
                    let list = await env.KV.get("ip", {
                        type: "json"
                    }) || [];
                    const body =await request.text();
                    await env.KV.put('ip',JSON.stringify(list))
                    return new Response(list,{
                        status: 200,
                        headers: {
                            'Content-Type': 'application/json;charset=utf-8'
                        }
                    })
                }
            }else if (urlrouter === 'deleteip'){
                await env.KV.put('ip',JSON.stringify([]))
                return new Response('ok',{
                    status: 200,
                })
            }
        }
        return new Response("hello world",{
            status: 200,
            headers: {
                'Content-Type': 'application/json;charset=utf-8'
            }
        })
    }
};
///////////////////////////////////////////////////////////////////////叉HTTP传输数据///////////////////////////////////////////////


function getCookie(cookie,name) {
    const cookies = cookie.split("; ");
    for (const cookie of cookies) {
        const [key, value] = cookie.split("=");
        if (key === name) {
            return decodeURIComponent(value);
        }
    }
    return null;
}
async function MD5MD5(文本) {
    const 编码器 = new TextEncoder();

    const 第一次哈希 = await crypto.subtle.digest('MD5', 编码器.encode(文本));
    const 第一次哈希数组 = Array.from(new Uint8Array(第一次哈希));
    const 第一次十六进制 = 第一次哈希数组.map(字节 => 字节.toString(16).padStart(2, '0')).join('');

    const 第二次哈希 = await crypto.subtle.digest('MD5', 编码器.encode(第一次十六进制.slice(7, 27)));
    const 第二次哈希数组 = Array.from(new Uint8Array(第二次哈希));
    const 第二次十六进制 = 第二次哈希数组.map(字节 => 字节.toString(16).padStart(2, '0')).join('');

    return 第二次十六进制.toLowerCase();
}