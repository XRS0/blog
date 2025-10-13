import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeRaw from 'rehype-raw';

interface MarkdownRendererProps {
  content: string;
}

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
  return (
    <div className="markdown-content">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeRaw]}
        components={{
          // Кастомные компоненты для различных элементов
          h1: ({ node, ...props }) => <h1 className="md-h1" {...props} />,
          h2: ({ node, ...props }) => <h2 className="md-h2" {...props} />,
          h3: ({ node, ...props }) => <h3 className="md-h3" {...props} />,
          h4: ({ node, ...props }) => <h4 className="md-h4" {...props} />,
          p: ({ node, ...props }) => <p className="md-p" {...props} />,
          a: ({ node, ...props }) => (
            <a className="md-link" target="_blank" rel="noopener noreferrer" {...props} />
          ),
          code: ({ node, className, children, ...props }: any) => {
            const inline = !className;
            return inline ? (
              <code className="md-inline-code" {...props}>
                {children}
              </code>
            ) : (
              <code className={className} {...props}>
                {children}
              </code>
            );
          },
          pre: ({ node, ...props }) => <pre className="md-pre" {...props} />,
          blockquote: ({ node, ...props }) => <blockquote className="md-blockquote" {...props} />,
          ul: ({ node, ...props }) => <ul className="md-ul" {...props} />,
          ol: ({ node, ...props }) => <ol className="md-ol" {...props} />,
          li: ({ node, ...props }) => <li className="md-li" {...props} />,
          table: ({ node, ...props }) => (
            <div className="md-table-wrapper">
              <table className="md-table" {...props} />
            </div>
          ),
          thead: ({ node, ...props }) => <thead className="md-thead" {...props} />,
          tbody: ({ node, ...props }) => <tbody className="md-tbody" {...props} />,
          tr: ({ node, ...props }) => <tr className="md-tr" {...props} />,
          th: ({ node, ...props }) => <th className="md-th" {...props} />,
          td: ({ node, ...props }) => <td className="md-td" {...props} />,
          hr: ({ node, ...props }) => <hr className="md-hr" {...props} />,
          img: ({ node, ...props }) => <img className="md-img" {...props} />,
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
