document.addEventListener('DOMContentLoaded', function() {
  const codeBlocks = document.querySelectorAll('pre > code');

  codeBlocks.forEach(function(codeBlock) {
    const pre = codeBlock.parentElement;

    // Count lines in the code block
    const text = codeBlock.textContent;
    const lines = text.trim().split('\n');
    const lineCount = lines.length;

    // Only enhance blocks with more than 20 lines
    if (lineCount <= 20) return;

    // Skip blocks explicitly marked as non-collapsible
    if (pre.hasAttribute('data-no-collapse')) return;

    // Create wrapper div
    const wrapper = document.createElement('div');
    wrapper.className = 'code-block-collapsible collapsed';

    // Create toggle button
    const button = document.createElement('button');
    button.className = 'code-toggle';
    button.setAttribute('aria-expanded', 'false');
    button.setAttribute('aria-label', 'Expand code block');
    button.textContent = 'Show all ' + lineCount + ' lines';

    // Wrap the pre element
    pre.parentNode.insertBefore(wrapper, pre);
    wrapper.appendChild(pre);
    wrapper.appendChild(button);

    // Toggle handler
    button.addEventListener('click', function() {
      const isCollapsed = wrapper.classList.contains('collapsed');

      if (isCollapsed) {
        wrapper.classList.remove('collapsed');
        button.setAttribute('aria-expanded', 'true');
        button.setAttribute('aria-label', 'Collapse code block');
        button.textContent = 'Collapse';
      } else {
        wrapper.classList.add('collapsed');
        button.setAttribute('aria-expanded', 'false');
        button.setAttribute('aria-label', 'Expand code block');
        button.textContent = 'Show all ' + lineCount + ' lines';
      }
    });
  });
});
