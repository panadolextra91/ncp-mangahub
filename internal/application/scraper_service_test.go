package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScraperService_ParseHTML(t *testing.T) {
	html := `
	<div class="quote">
        <span class="text">“The world as we have created it is a process of our thinking. It cannot be changed without changing our thinking.”</span>
        <span>by <small class="author">Albert Einstein</small></span>
        <div class="tags">
            <a class="tag">change</a>
            <a class="tag">deep-thoughts</a>
        </div>
    </div>
	<div class="quote">
        <span class="text">“It is our choices, Harry, that show what we truly are, far more than our abilities.”</span>
        <span>by <small class="author">J.K. Rowling</small></span>
        <div class="tags">
            <a class="tag">abilities</a>
            <a class="tag">choices</a>
        </div>
    </div>
	`

	s := NewScraperService()
	quotes := s.parseHTML(html)

	assert.Len(t, quotes, 2)
	
	assert.Equal(t, "The world as we have created it is a process of our thinking. It cannot be changed without changing our thinking.", quotes[0].Text)
	assert.Equal(t, "Albert Einstein", quotes[0].Author)
	assert.Contains(t, quotes[0].Tags, "change")
	assert.Contains(t, quotes[0].Tags, "deep-thoughts")

	assert.Equal(t, "It is our choices, Harry, that show what we truly are, far more than our abilities.", quotes[1].Text)
	assert.Equal(t, "J.K. Rowling", quotes[1].Author)
	assert.Contains(t, quotes[1].Tags, "choices")
}

func TestScraperService_CleanText(t *testing.T) {
	s := NewScraperService()
	assert.Equal(t, "“Hello”", s.cleanText("&ldquo;Hello&rdquo;"))
	assert.Equal(t, "It's me", s.cleanText("It&rsquo;s me"))
}
