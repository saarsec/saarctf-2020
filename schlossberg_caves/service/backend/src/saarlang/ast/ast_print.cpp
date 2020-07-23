#include "ast.h"
#include "statements.h"
#include "expressions.h"


using namespace std;

/*
 * For debugging: prettyprint a given AST node.
 * AST => sourcecode
 */


void ConstantDefNode::print(std::ostream &out) {
	out << "const " << symbol << ": ";
	dynamic_cast<TypeNode *>(children[0])->print(out);
	out << " = ";
	dynamic_cast<ExpressionNode *>(children[1])->print(out);
	out << ";\n";
}

void FunctionDefNode::print(std::ostream &out) {
	out << "eija " << symbol << " (";
	bool first = true;
	for (auto &it: arguments) {
		if (first) first = false;
		else out << ", ";
		out << it->symbol << ": ";
		it->getType()->print(out);
	}
	out << ") gebbtserick ";
	dynamic_cast<TypeNode *>(children[0])->print(out);
	out << ": \n";
	dynamic_cast<StatementNode *>(children[1])->print(out, 0);
}

void ImportNode::print(std::ostream &out) {
	out << "holmol \"" << filename << "\";\n";
}

void ExpressionStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	dynamic_cast<ExpressionNode *>(children[0])->print(out);
	out << ";\n";
}

void IfStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "falls ";
	dynamic_cast<ExpressionNode *>(children[0])->print(out);
	out << ":\n";
	dynamic_cast<StatementNode *>(children[1])->print(out, indent + 1);
	if (children.size() > 2) {
		indentation(out, indent);
		out << "sonschd:\n";
		dynamic_cast<StatementNode *>(children[2])->print(out, indent + 1);
	}
}

void WhileStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "solang ";
	dynamic_cast<ExpressionNode *>(children[0])->print(out);
	out << ":\n";
	dynamic_cast<StatementNode *>(children[1])->print(out, indent + 1);
}

void VarDefStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "var " << symbol << ": ";
	dynamic_cast<TypeNode *>(children[0])->print(out);
	if (children.size() > 1) {
		out << " = ";
		dynamic_cast<ExpressionNode *>(children[1])->print(out);
	}
	out << ";\n";
}

void VarArrayStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "var " << symbol << ": ";
	dynamic_cast<TypeNode *>(children[0])->print(out);
	out << " (" << size << ");\n";
}

void ReturnStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "serick";
	if (!children.empty()) {
		out << " ";
		dynamic_cast<ExpressionNode *>(children[0])->print(out);
	}
	out << ";\n";
}

void BlockStmtNode::print(std::ostream &out, int indent) {
	indentation(out, indent);
	out << "{\n";
	for (const auto &it: children) {
		dynamic_cast<StatementNode *>(it)->print(out, indent + 1);
	}
	indentation(out, indent);
	out << "}\n";
}

void SymbolExprNode::print(std::ostream &out) {
	out << token.text;
}

void ConstantExprNode::print(std::ostream &out) {
	out << token.text;
}

void CallExprNode::print(std::ostream &out) {
	out << "mach ";
	dynamic_cast<ExpressionNode *>(children[0])->print(out);
	out << " (";
	for (int i = 1; i < children.size(); i++) {
		dynamic_cast<ExpressionNode *>(children[i])->print(out);
		if (i < children.size() - 1) out << ", ";
	}
	out << ")";
}

void NewArrayExprNode::print(std::ostream &out) {
	out << "neie ";
	dynamic_cast<TypeNode *>(children[0])->print(out);
	out << " (";
	dynamic_cast<ExpressionNode *>(children[1])->print(out);
	out << ")";
}

void ArrayLengthExprNode::print(std::ostream &out) {
	out << "grees ";
	dynamic_cast<ExpressionNode *>(children[0])->print(out);
}

void BinaryOperatorExprNode::print(std::ostream &out) {
	if (dynamic_cast<BinaryOperatorExprNode *>(children[0])) {
		out << "(";
		dynamic_cast<ExpressionNode *>(children[0])->print(out);
		out << ")";
	} else {
		dynamic_cast<ExpressionNode *>(children[0])->print(out);
	}

	out << " " << token.text << " ";

	if (dynamic_cast<BinaryOperatorExprNode *>(children[1])) {
		out << "(";
		dynamic_cast<ExpressionNode *>(children[1])->print(out);
		out << ")";
	} else {
		dynamic_cast<ExpressionNode *>(children[1])->print(out);
	}
}

void IntrinsicExprNode::print(std::ostream &out) {
	out << token.text;
}