#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>
#include "sandbox.c"

#define HF_MAX_CONTENT_SZ 10*1096

#if (NGX_DEBUG)
#define HT_SAARJSF_DEBUG 0
#else
#define HT_SAARJSF_DEBUG 1
#endif

/*
stack for parsing html
*/
typedef struct 
{
u_char data[HF_MAX_CONTENT_SZ];
ngx_int_t top;
}
saarjsfilter_stack_t;


/*
module data struct for maintaining
state per request
*/
typedef struct
{
ngx_uint_t  last;
ngx_uint_t  count;
ngx_uint_t  found;
ngx_uint_t  index;
ngx_uint_t  starttag;
ngx_uint_t  inner_length;  
ngx_uint_t  outer_offset;  
ngx_uint_t  tagquote;
ngx_uint_t  tagsquote;
saarjsfilter_stack_t stack;
ngx_chain_t  *free;
ngx_chain_t  *busy;
ngx_chain_t  *out;
ngx_chain_t  *in;
ngx_chain_t  **last_out;
}
ngx_http_html_saarjs_filter_ctx_t;



static ngx_http_output_header_filter_pt  ngx_http_next_header_filter;
static ngx_http_output_body_filter_pt    ngx_http_next_body_filter;


/* Function prototypes */
static ngx_int_t ngx_http_html_saarjs_init(ngx_conf_t * cf);
static ngx_int_t ngx_http_html_saarjs_header_filter(ngx_http_request_t *r );
static ngx_int_t ngx_http_html_saarjs_body_filter(ngx_http_request_t *r, 
                                                ngx_chain_t *in);
static ngx_int_t ngx_test_content_type(ngx_http_request_t *r);
static ngx_int_t ngx_test_content_compression(ngx_http_request_t *r);
static void ngx_init_stack(saarjsfilter_stack_t *stack);
static ngx_int_t push(u_char c, saarjsfilter_stack_t *stack);
static ngx_int_t pop(saarjsfilter_stack_t *stack);


//static ngx_int_t ngx_process_tag(ngx_http_html_saarjs_filter_ctx_t *ctx, 
//                                 ngx_http_request_t *r);

static ngx_int_t ngx_html_insert_output(
                    ngx_http_html_saarjs_filter_ctx_t *ctx, 
                    ngx_http_request_t *r, ngx_str_t *str);

                    


/*
Module context 
*/
static ngx_http_module_t  ngx_http_html_saarjs_filter_ctx =
{
    NULL, //Pre config
    ngx_http_html_saarjs_init, //Post config
    NULL, //Create main config
    NULL, //Init main config
    NULL, //Create server config
    NULL, //Merge server config
    NULL, //Create loc config
    NULL //Merge loc config
};


/*
Module definition
*/
ngx_module_t  ngx_http_html_saarjs_filter_module = 
{
    NGX_MODULE_V1,
    &ngx_http_html_saarjs_filter_ctx,     /* module context */
    NULL, /* module directives */
    NGX_HTTP_MODULE,                    /* module type */
    NULL,                                  
    NULL,                                  
    NULL,                                  
    NULL,                                  
    NULL,                                  
    NULL,                                 
    NULL,                                  
    NGX_MODULE_V1_PADDING
};



/* Function to initialize the module */
static ngx_int_t
ngx_http_html_saarjs_init(ngx_conf_t * cfg)
{

    ngx_http_next_header_filter = ngx_http_top_header_filter;
    ngx_http_top_header_filter = ngx_http_html_saarjs_header_filter;

    ngx_http_next_body_filter = ngx_http_top_body_filter;
    ngx_http_top_body_filter = ngx_http_html_saarjs_body_filter;

 
    return NGX_OK;

}


/*Module function handler to filter http response headers */
static ngx_int_t
ngx_http_html_saarjs_header_filter(ngx_http_request_t *r )
{

    ngx_http_html_saarjs_filter_ctx_t *ctx;
    //ngx_uint_t content_length=0; 
    

    if(r->headers_out.content_type.len == 0 || 
        r->headers_out.content_length_n == 0 ||
        r->header_only )
    {
        #if HT_SAARJSF_DEBUG
            ngx_log_debug0(NGX_LOG_DEBUG_HTTP, r->connection->log, 0,
                "[Html_saarjs filter]: empty content type or "
                "header only ");
        #endif
        
        return ngx_http_next_header_filter(r);
    }
    
     
    if(ngx_test_content_type(r) == 0) 
    {
        #if HT_SAARJSF_DEBUG
            ngx_log_debug0(NGX_LOG_DEBUG_HTTP, r->connection->log, 0,
                "[Html_saarjs filter]: content type not html");
        #endif            
        
        return ngx_http_next_header_filter(r);
    }

    
    if(ngx_test_content_compression(r) != 0)
    {//Compression enabled, don't filter   
        ngx_log_stderr( 0, 
                     "[Html_saarjs filter]: compression enabled");
                     
        return ngx_http_next_header_filter(r);
    }
 
    if(r->headers_out.status != NGX_HTTP_OK)
    {//Response is not HTTP 200   
        //ngx_log_stderr( 0, 
        //             "[Html_saarjs filter]: http response is not 200");
                     
        return ngx_http_next_header_filter(r);
    }

    r->filter_need_in_memory = 1;

    if (r == r->main) 
    {//Main request 
        ngx_http_clear_content_length(r);  
    }
    

    ctx = ngx_http_get_module_ctx(r, ngx_http_html_saarjs_filter_module);
    if(ctx == NULL)
    {
        ctx = ngx_pcalloc(r->pool, 
                          sizeof(ngx_http_html_saarjs_filter_ctx_t)); 
        
        if(ctx == NULL)
        {
            ngx_log_stderr( 0,
                          "[Html_saarjs filter]: cannot allocate ctx"
                          " memory");
                          
            return ngx_http_next_header_filter(r);
        }
        
        ngx_http_set_ctx(r, ctx, ngx_http_html_saarjs_filter_module);
    }
    
    
    return ngx_http_next_header_filter(r);
    
}



/*
Module function handler to filter the html response body
and insert the text string 
*/
static ngx_int_t
ngx_http_html_saarjs_body_filter(ngx_http_request_t *r, ngx_chain_t *in)
{

    ngx_http_html_saarjs_filter_ctx_t *ctx;
    //ngx_chain_t  *cl;
    //ngx_buf_t  *b;
    ngx_int_t rc;
                                  

    ctx = ngx_http_get_module_ctx(r, ngx_http_html_saarjs_filter_module);


    if(ctx == NULL)
    {
        #if HT_SAARJSF_DEBUG
            ngx_log_debug0(NGX_LOG_DEBUG_HTTP, r->connection->log, 0,
                "[Html_saarjs filter]: ngx_http_html_saarjs_body_filter "
                "unable to get module ctx");
        #endif           
            
        return ngx_http_next_body_filter(r, in);
    }


    if(in == NULL)
    {
       ngx_log_stderr( 0, 
            "[Html_saarjs filter]: input chain is null");
                     
       return ngx_http_next_body_filter(r, in);
    }


    //Copy the incoming chain to ctx-in
    if (ngx_chain_add_copy(r->pool, &ctx->in, in) != NGX_OK) 
    {
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: unable to copy"
            " input chain - in");
                     
        return NGX_ERROR;
    }

    ctx->last_out = &ctx->out;
   
    uint tagid = 0;
    uint tagid2 = 1;
    //while (tagid == 0 || tagid != tagid2){
    //Loop through all the incoming buffers
    while(ctx->in)
    {	
        ctx->index = 0; 
        u_char *p, c;
        u_char pbuf[12]; 
        //ngx_int_t rc;
        ngx_buf_t* buf;
        
        if(ctx->in == NULL)
        {
            ngx_log_stderr( 0, 
                "[Html_saarjs filter]: ngx_parse_buf_html "
                "unable to parse html ctx->in is NULL");  
                
            return NGX_ERROR;
        }
            
            buf = ctx->in->buf; 
            ngx_memset(pbuf, 0, sizeof(pbuf));
            uint inner_scope = 0;
            for(p=buf->pos; p < buf->last; p++)
            {

                c = *p;
                pbuf[sizeof(pbuf)-1] = 0;
                for (uint i = 1; i < sizeof(pbuf) - 1; ++i)
                    pbuf[i - 1] = pbuf[i];
                pbuf[sizeof(pbuf)- 2] = c;
                //ngx_log_stderr( 0, 
                
                //    (char *)pbuf);  
                ctx->outer_offset++;
                if (!inner_scope){
                    ctx->index++;
                    if (ngx_strncmp(pbuf+(sizeof(pbuf)-9), "<?saarjs", 8) == 0){
                        tagid++;
                        if (tagid == tagid2){ 
                            tagid2++;
                            ctx->index -= 9;
                            inner_scope = 1;
                            ngx_init_stack(&ctx->stack);
                        }
                    }
                }
                else {
                    if (ngx_strncmp(pbuf+(sizeof(pbuf)-3), "?>", 2) == 0){
                        inner_scope = 0;
                        pop(&ctx->stack);
                        pop(&ctx->stack);
                        struct timespec max_wait;
                        memset(&max_wait, 0, sizeof(max_wait));
                        /* wait at most 2 seconds */
                        max_wait.tv_sec = 2;
                        do_or_timeout(&max_wait, (char *) ctx->stack.data);
                        ngx_str_t *result = malloc(sizeof(ngx_str_t));
                        result->data = malloc(sizeof(saarjs_buf));
                        ngx_memcpy(result->data, saarjs_buf, sizeof(saarjs_buf));
                        result->len = ngx_strlen(result->data);
                        // TODO: execute js
                        ngx_html_insert_output(ctx, r, result);
                        free(result);

                        break;
                    } else {
                        if (push(c, &ctx->stack)){
                            ngx_log_stderr( 0, 
                                "Script exceeds limit");  
                                
                            return NGX_ERROR;
                        }
                    }
                }
            }
            if (inner_scope != 0) {
                continue;
            }
        
        //b = ctx->in->buf;	
    
        *ctx->last_out=ctx->in;
        ctx->last_out=&ctx->in->next;
        ctx->in = ctx->in->next;
    } 
    //}
    //ngx_log_stderr( 0, 
    //    "end while");  

    *ctx->last_out = NULL;
    
   
    rc=ngx_http_next_body_filter(r, ctx->out);

    ngx_chain_update_chains(r->pool, &ctx->free, &ctx->busy, &ctx->out,
                            (ngx_buf_tag_t)&ngx_http_html_saarjs_filter_module);

    ctx->in = NULL; 

    return rc;
    
}


/*
Insert the text into body response buffer
*/
static ngx_int_t 
ngx_html_insert_output(ngx_http_html_saarjs_filter_ctx_t *ctx, 
                       ngx_http_request_t *r, 
                       ngx_str_t *slcf)
{

    ngx_chain_t  *cl, *ctx_in_new, **ll;
    ngx_buf_t  *b;

    if(ctx->in == NULL)
    {
        ngx_log_stderr( 0, 
             "[Html_saarjs filter]: ngx_html_insert_output "
             "text Insertion ctx->in is NULL");
             
        return NGX_ERROR;
    }

				   
    ll = &ctx_in_new;				   
    b=ctx->in->buf;

    if(b->pos + ctx->index + 1 > b->last)
    {//Check that the saarjs tag position does not exceed buffer
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: ngx_html_insert_output "
            "invalid input buffer at text insertion");
            
        return NGX_ERROR;          
    }

    cl = ngx_chain_get_free_buf(r->pool, &ctx->free);
    if (cl == NULL) 
    {
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: ngx_html_insert_output "
            "unable to allocate output chain");
            
        return NGX_ERROR;
    }

    b=cl->buf;
    ngx_memzero(b, sizeof(ngx_buf_t));

    b->tag = (ngx_buf_tag_t) &ngx_http_html_saarjs_filter_module;
    b->memory=1;
    b->pos = ctx->in->buf->pos;
    b->last = b->pos + ctx->index + 1;
    b->recycled = ctx->in->buf->recycled;
    b->flush = ctx->in->buf->flush; 
       
    *ll = cl;  
    ll = &cl->next;
	

    cl = ngx_chain_get_free_buf(r->pool, &ctx->free);
    if (cl == NULL) 
    {
        ngx_log_stderr( 0, 
             "[Html_saarjs filter]: ngx_html_insert_output "
             "unable to allocate output chain");
             
        return NGX_ERROR;
    }

    b=cl->buf;
    ngx_memzero(b, sizeof(ngx_buf_t));
	 
    b->tag = (ngx_buf_tag_t) &ngx_http_html_saarjs_filter_module;
    b->memory=1;
    b->pos=slcf->data;
    b->last=b->pos + slcf->len;
    b->recycled = ctx->in->buf->recycled;
	 
    *ll = cl;
    ll = &cl->next;
	 

    if(ctx->in->buf->pos + ctx->index + 1 == ctx->in->buf->last )
    {//saarjs tag is in last position of the buffer
   
        b->last_buf = ctx->in->buf->last_buf;
        b->last_in_chain = ctx->in->buf->last_in_chain;
		 
        *ll = ctx->in->next; 
		
	    if(ctx->in->buf->recycled)
	    {//consume existing buffer
	        ctx->in->buf->pos = ctx->in->buf->last;
	    }
	    ctx->in = ctx_in_new;
	    return NGX_OK;
		
    }
     
    
    //tag is within buffer last position, 
    //i.e. ctx->in->buf->pos + ctx->index + 1 < ctx->in->buf->last
    cl = ngx_chain_get_free_buf(r->pool, &ctx->free);
    if (cl == NULL) 
    {
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: ngx_html_insert_output unable to allocate "
            "output chain");
            
        return NGX_ERROR;
    }

    b=cl->buf;
    ngx_memzero(b, sizeof(ngx_buf_t));

    b->tag = (ngx_buf_tag_t) &ngx_http_html_saarjs_filter_module;
    b->memory=1;
    ctx->index += ctx->stack.top + 13;
    b->pos = ctx->in->buf->pos + ctx->index;
    b->last = ctx->in->buf->last;
    b->recycled = ctx->in->buf->recycled;
    b->last_buf = ctx->in->buf->last_buf;
    b->last_in_chain = ctx->in->buf->last_in_chain;

    *ll = cl;
    ll = &cl->next;
    *ll = ctx->in->next;
	 
    if(ctx->in->buf->recycled)
    {//consume existing buffer
        ctx->in->buf->pos = ctx->in->buf->last;	
    }
	 
    ctx->in = ctx_in_new; 

    
	   
    return NGX_OK;

}


/*
Check if the content is text/html 
Returns true if text/html is present, false otherwise
*/
static ngx_int_t
ngx_test_content_type(ngx_http_request_t *r)
{
    ngx_str_t tmp;

    if(r->headers_out.content_type.len == 0)
    {
        return 0;
    } 

    tmp.len = r->headers_out.content_type.len;
    tmp.data = ngx_pcalloc(r->pool, sizeof(u_char) * tmp.len ); 

    if(tmp.data == NULL)
    {
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: ngx_test_content_type "
            "cannot allocate buffer for content type check");
        return 0;
    } 

    ngx_strlow(tmp.data, r->headers_out.content_type.data, tmp.len); 

    if(ngx_strnstr(tmp.data, "text/html", 
                  r->headers_out.content_type.len) != NULL)
    {
        return 1;
    }
   
    return 0; 
    
}


/*
Check if the content encoding is compressed using either
gzip, deflate, compress or br (Brotli)
Returns true if compression is enabled, 
false if it cannot determine compression
*/
static ngx_int_t
ngx_test_content_compression(ngx_http_request_t *r)
{
    ngx_str_t tmp;
    
    if(r->headers_out.content_encoding == NULL)
    {//Cannot determine encoding, assume no compression
        return 0; 
    }

    if(r->headers_out.content_encoding->value.len == 0 )
    {
        return 0; 
    }

    tmp.len = r->headers_out.content_encoding->value.len;
    tmp.data = ngx_pcalloc(r->pool, sizeof(u_char) * tmp.len );

    if(tmp.data == NULL)
    {
        ngx_log_stderr( 0, 
            "[Html_saarjs filter]: ngx_test_content_compression"
            " cannot allocate buffer for compression check");
            
        return 0;
    }

    ngx_strlow(tmp.data, 
               r->headers_out.content_encoding->value.data, tmp.len); 


    
    if( tmp.len >= (sizeof("gzip") -1) && 
        ngx_strncmp(tmp.data, (u_char*)"gzip" , tmp.len) == 0 )
    {
        return 1; 
    }
    
    if( tmp.len >= (sizeof("deflate") -1) &&
        ngx_strncmp(tmp.data, (u_char*)"deflate" , tmp.len) == 0 )
    {
        return 1; 
    }
    
    if( tmp.len >= (sizeof("compress") -1) &&
        ngx_strncmp(tmp.data, (u_char*)"compress" , tmp.len) == 0 )
    {
        return 1; 
    }
    
   
    if( tmp.len >= (sizeof("br") -1) &&
        ngx_strncmp(tmp.data, (u_char*)"br" , tmp.len) == 0 )
    {
        return 1; 
    }
        
    //Fail safe to false if compression cannot be determined
    return 0; 
}

/*
Initializes the stack structure
*/
static void 
ngx_init_stack(saarjsfilter_stack_t *stack)
{
    ngx_memset(stack, 0 , sizeof(saarjsfilter_stack_t)); 
    stack->top = -1; 
}

/*
Push a u_char into the stack 
Returns -1 if out of stack space 
*/
static ngx_int_t 
push(u_char c, saarjsfilter_stack_t *stack)
{

    if(stack->top == (HF_MAX_CONTENT_SZ -1) )
       return -1;
    
    stack->top++;
    stack->data[stack->top] = c;
    return 0;    
}

static ngx_int_t 
pop(saarjsfilter_stack_t *stack)
{
 
    if(stack->top == 0 )
       return 0;
    
    ngx_int_t c = stack->data[stack->top];
    stack->data[stack->top] = 0;
    stack->top--;
    return c;    
}
