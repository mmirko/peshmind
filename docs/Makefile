TEX_IT=relazione.tex
TEX_EN=relazione_en.tex
PDF_IT=relazione.pdf
PDF_EN=relazione_en.pdf

.PHONY: all pdf pdf-it pdf-en clean

all: pdf

pdf: pdf-it pdf-en

pdf-it:
	latexmk -pdf -interaction=nonstopmode -halt-on-error $(TEX_IT)

pdf-en:
	latexmk -pdf -interaction=nonstopmode -halt-on-error $(TEX_EN)

clean:
	latexmk -C
	rm -f *.vrb *.nav *.snm $(PDF_IT) $(PDF_EN)
